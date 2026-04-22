package importer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kest-labs/kest/api/internal/modules/collection"
	"github.com/kest-labs/kest/api/internal/modules/request"
)

var (
	ErrInvalidMarkdownDocument = errors.New("invalid markdown API document")
	ErrMarkdownBaseURLNotFound = errors.New("unable to derive absolute request URLs from markdown document")
	ErrNoImportableEndpoints   = errors.New("no importable endpoints found in markdown document")
)

type MarkdownImportResult struct {
	RootFolderID       uint                         `json:"root_folder_id"`
	RootFolderName     string                       `json:"root_folder_name"`
	CollectionsCreated int                          `json:"collections_created"`
	RequestsCreated    int                          `json:"requests_created"`
	Modules            []MarkdownImportModuleResult `json:"modules"`
}

type MarkdownImportModuleResult struct {
	Name         string `json:"name"`
	CollectionID uint   `json:"collection_id"`
	RequestCount int    `json:"request_count"`
}

type markdownDocument struct {
	Title   string
	BaseURL string
	Modules []markdownModule
}

type markdownModule struct {
	Name      string
	Endpoints []markdownEndpoint
}

type markdownEndpoint struct {
	Name        string
	Description string
	Method      string
	Path        string
	URL         string
	Headers     []request.KeyValue
	QueryParams []request.KeyValue
	PathParams  map[string]string
	Body        string
	BodyType    string
}

type markdownSection struct {
	Title string
	Body  string
}

type curlExample struct {
	URL     string
	Headers []request.KeyValue
	Body    string
}

var (
	endpointHeadingPattern = regexp.MustCompile("^([A-Z]+)\\s+`([^`]+)`$")
	urlPattern             = regexp.MustCompile("https?://[^\\s'\"\\\\`]+")
	curlHeaderPattern      = regexp.MustCompile(`-H\s+(?:'([^']*)'|"([^"]*)")`)
	curlDataPattern        = regexp.MustCompile(`-d\s+(?:'([^']*)'|"([^"]*)")`)
	boldLinePattern        = regexp.MustCompile(`(?m)^\*\*(.+?)\*\*$`)
)

func (s *service) ImportMarkdown(
	ctx context.Context,
	projectID, parentID uint,
	file *multipart.FileHeader,
) (*MarkdownImportResult, error) {
	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	doc, err := parseMarkdownDocument(file.Filename, content)
	if err != nil {
		return nil, err
	}

	return s.importMarkdownDocument(ctx, projectID, parentID, doc)
}

func (s *service) importMarkdownDocument(
	ctx context.Context,
	projectID, parentID uint,
	doc *markdownDocument,
) (*MarkdownImportResult, error) {
	if doc == nil || len(doc.Modules) == 0 {
		return nil, ErrNoImportableEndpoints
	}

	rootReq := &collection.CreateCollectionRequest{
		ProjectID: projectID,
		Name:      doc.Title,
		IsFolder:  true,
	}
	if parentID > 0 {
		rootReq.ParentID = &parentID
	}

	rootFolder, err := s.collectionService.Create(ctx, rootReq)
	if err != nil {
		return nil, err
	}

	result := &MarkdownImportResult{
		RootFolderID:   rootFolder.ID,
		RootFolderName: rootFolder.Name,
		Modules:        make([]MarkdownImportModuleResult, 0, len(doc.Modules)),
	}

	for moduleIndex, module := range doc.Modules {
		moduleReq := &collection.CreateCollectionRequest{
			ProjectID: projectID,
			Name:      module.Name,
			ParentID:  &rootFolder.ID,
			IsFolder:  false,
			SortOrder: moduleIndex,
		}

		moduleCollection, err := s.collectionService.Create(ctx, moduleReq)
		if err != nil {
			return nil, err
		}

		moduleResult := MarkdownImportModuleResult{
			Name:         module.Name,
			CollectionID: moduleCollection.ID,
		}

		for requestIndex, endpoint := range module.Endpoints {
			req := &request.CreateRequestRequest{
				CollectionID: moduleCollection.ID,
				Name:         endpoint.Name,
				Description:  endpoint.Description,
				Method:       endpoint.Method,
				URL:          endpoint.URL,
				Headers:      endpoint.Headers,
				QueryParams:  endpoint.QueryParams,
				PathParams:   endpoint.PathParams,
				Body:         endpoint.Body,
				BodyType:     endpoint.BodyType,
				SortOrder:    requestIndex,
			}

			if _, err := s.requestService.Create(ctx, projectID, req); err != nil {
				return nil, err
			}

			moduleResult.RequestCount++
			result.RequestsCreated++
		}

		result.CollectionsCreated++
		result.Modules = append(result.Modules, moduleResult)
	}

	return result, nil
}

func parseMarkdownDocument(filename string, content []byte) (*markdownDocument, error) {
	text := normalizeMarkdown(string(content))
	if text == "" {
		return nil, ErrInvalidMarkdownDocument
	}

	title := extractDocumentTitle(text)
	if title == "" {
		title = fallbackDocumentTitle(filename)
	}

	doc := &markdownDocument{
		Title:   title,
		BaseURL: extractDocumentBaseURL(text),
	}

	if looksLikeSingleModuleDocument(text) {
		endpoints, err := parseEndpointSections(text, doc.BaseURL)
		if err != nil {
			return nil, err
		}
		if len(endpoints) == 0 {
			return nil, ErrNoImportableEndpoints
		}
		doc.Modules = append(doc.Modules, markdownModule{
			Name:      normalizeSingleModuleName(title),
			Endpoints: endpoints,
		})
		return doc, nil
	}

	for _, section := range splitByHeadingLevel(text, 2) {
		endpoints, err := parseEndpointSections(section.Body, doc.BaseURL)
		if err != nil {
			return nil, err
		}
		if len(endpoints) == 0 {
			continue
		}

		doc.Modules = append(doc.Modules, markdownModule{
			Name:      strings.TrimSpace(section.Title),
			Endpoints: endpoints,
		})
	}

	if len(doc.Modules) == 0 {
		return nil, ErrNoImportableEndpoints
	}

	return doc, nil
}

func parseEndpointSections(content, baseURL string) ([]markdownEndpoint, error) {
	sections := splitByHeadingLevel(content, 3)
	endpoints := make([]markdownEndpoint, 0, len(sections))

	for _, section := range sections {
		matches := endpointHeadingPattern.FindStringSubmatch(strings.TrimSpace(section.Title))
		if len(matches) != 3 {
			continue
		}

		endpoint, err := parseEndpointSection(matches[1], matches[2], section.Body, baseURL)
		if err != nil {
			return nil, err
		}

		endpoints = append(endpoints, endpoint)
	}

	return endpoints, nil
}

func parseEndpointSection(method, path, body, baseURL string) (markdownEndpoint, error) {
	endpoint := markdownEndpoint{
		Method: strings.ToUpper(strings.TrimSpace(method)),
		Path:   strings.TrimSpace(path),
	}

	summary := extractBoldSummary(body)
	if summary == "" {
		summary = fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path)
	}
	endpoint.Name = summary
	endpoint.Description = summary

	subsections := splitByHeadingLevel(body, 4)
	sectionMap := make(map[string]string, len(subsections))
	for _, section := range subsections {
		sectionMap[strings.ToLower(strings.TrimSpace(section.Title))] = section.Body
	}

	pathDefaults := parseParameterDefaults(sectionMap["path parameters"])
	example := parseCurlExample(sectionMap["example"])

	finalURL, err := buildAbsoluteEndpointURL(baseURL, example.URL, endpoint.Path)
	if err != nil {
		return markdownEndpoint{}, err
	}
	endpoint.URL = finalURL

	endpoint.PathParams = pathDefaults
	for key, value := range extractPathParamsFromExample(endpoint.Path, example.URL) {
		endpoint.PathParams[key] = value
	}

	endpoint.QueryParams = parseQueryParamsFromExample(example.URL)
	endpoint.Headers = filterNonAuthHeaders(example.Headers)
	if strings.Contains(body, "JWT Required") {
		endpoint.Headers = append(endpoint.Headers, request.KeyValue{
			Key:     "Authorization",
			Value:   "Bearer <token>",
			Enabled: false,
		})
	}

	requestBodyLang, requestBody := firstCodeBlock(sectionMap["request body"])
	switch {
	case requestBody != "":
		endpoint.Body = requestBody
		endpoint.BodyType = deriveBodyType(requestBodyLang, endpoint.Headers)
	case example.Body != "":
		endpoint.Body = example.Body
		endpoint.BodyType = deriveBodyType("", endpoint.Headers)
	default:
		endpoint.BodyType = "none"
	}

	if len(endpoint.PathParams) == 0 {
		endpoint.PathParams = nil
	}

	return endpoint, nil
}

func normalizeMarkdown(content string) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.TrimPrefix(content, "\ufeff")
	return strings.TrimSpace(content)
}

func extractDocumentTitle(content string) string {
	for _, line := range strings.Split(content, "\n") {
		level, title, ok := parseHeading(line)
		if ok && level == 1 {
			return strings.TrimSpace(title)
		}
	}
	return ""
}

func fallbackDocumentTitle(filename string) string {
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	return strings.TrimSpace(strings.TrimSuffix(base, ext))
}

func normalizeSingleModuleName(title string) string {
	if trimmed := strings.TrimSpace(strings.TrimSuffix(title, " API")); trimmed != "" {
		return trimmed
	}
	return title
}

func extractDocumentBaseURL(content string) string {
	for _, section := range splitByHeadingLevel(content, 2) {
		title := strings.TrimSpace(section.Title)
		if title != "基础 URL" && title != "Base URL" {
			continue
		}

		if match := urlPattern.FindString(section.Body); match != "" {
			return strings.TrimSpace(match)
		}
	}
	return ""
}

func looksLikeSingleModuleDocument(content string) bool {
	return strings.Contains(content, "\n## Endpoints\n") || strings.Contains(content, "\n## Details\n")
}

func splitByHeadingLevel(content string, level int) []markdownSection {
	lines := strings.Split(content, "\n")
	sections := make([]markdownSection, 0)
	var current *markdownSection

	for _, line := range lines {
		lineLevel, title, ok := parseHeading(line)
		if ok && lineLevel == level {
			if current != nil {
				current.Body = strings.TrimSpace(current.Body)
				sections = append(sections, *current)
			}
			current = &markdownSection{Title: strings.TrimSpace(title)}
			continue
		}

		if current == nil {
			continue
		}

		if current.Body == "" {
			current.Body = line
		} else {
			current.Body += "\n" + line
		}
	}

	if current != nil {
		current.Body = strings.TrimSpace(current.Body)
		sections = append(sections, *current)
	}

	return sections
}

func parseHeading(line string) (int, string, bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "#") {
		return 0, "", false
	}

	level := 0
	for level < len(trimmed) && trimmed[level] == '#' {
		level++
	}
	if level == 0 || level >= len(trimmed) || trimmed[level] != ' ' {
		return 0, "", false
	}

	return level, strings.TrimSpace(trimmed[level+1:]), true
}

func extractBoldSummary(content string) string {
	match := boldLinePattern.FindStringSubmatch(content)
	if len(match) != 2 {
		return ""
	}
	return strings.TrimSpace(match[1])
}

func parseParameterDefaults(section string) map[string]string {
	rows := parseMarkdownTable(section)
	if len(rows) <= 1 {
		return nil
	}

	result := make(map[string]string)
	for _, row := range rows[1:] {
		if len(row) < 2 {
			continue
		}
		name := unwrapMarkdownLiteral(row[0])
		if name == "" {
			continue
		}
		result[name] = defaultValueForType(unwrapMarkdownLiteral(row[1]))
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func parseMarkdownTable(content string) [][]string {
	rows := make([][]string, 0)
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") {
			continue
		}

		cells := splitMarkdownTableRow(trimmed)
		if len(cells) == 0 || isMarkdownSeparatorRow(cells) {
			continue
		}
		rows = append(rows, cells)
	}
	return rows
}

func splitMarkdownTableRow(line string) []string {
	line = strings.TrimSpace(strings.Trim(line, "|"))
	if line == "" {
		return nil
	}

	rawCells := strings.Split(line, "|")
	cells := make([]string, 0, len(rawCells))
	for _, cell := range rawCells {
		cells = append(cells, strings.TrimSpace(cell))
	}
	return cells
}

func isMarkdownSeparatorRow(cells []string) bool {
	for _, cell := range cells {
		trimmed := strings.Trim(strings.TrimSpace(cell), ":")
		if trimmed == "" || strings.Trim(trimmed, "-") != "" {
			return false
		}
	}
	return true
}

func unwrapMarkdownLiteral(value string) string {
	return strings.Trim(strings.TrimSpace(value), "`")
}

func defaultValueForType(typeName string) string {
	normalized := strings.ToLower(strings.TrimSpace(typeName))
	switch {
	case strings.Contains(normalized, "int"), strings.Contains(normalized, "uint"), strings.Contains(normalized, "number"):
		return "1"
	case strings.Contains(normalized, "bool"):
		return "true"
	case strings.Contains(normalized, "float"), strings.Contains(normalized, "double"), strings.Contains(normalized, "decimal"):
		return "1"
	default:
		return "example"
	}
}

func parseCurlExample(section string) curlExample {
	if section == "" {
		return curlExample{}
	}

	_, code := firstCodeBlock(section)
	if code == "" {
		code = section
	}

	example := curlExample{}
	if match := urlPattern.FindString(code); match != "" {
		example.URL = strings.TrimSpace(match)
	}

	for _, match := range curlHeaderPattern.FindAllStringSubmatch(code, -1) {
		headerText := firstNonEmpty(match[1], match[2])
		if headerText == "" {
			continue
		}

		key, value, ok := strings.Cut(headerText, ":")
		if !ok {
			continue
		}

		example.Headers = append(example.Headers, request.KeyValue{
			Key:     strings.TrimSpace(key),
			Value:   strings.TrimSpace(value),
			Enabled: true,
		})
	}

	if match := curlDataPattern.FindStringSubmatch(code); len(match) > 0 {
		example.Body = strings.TrimSpace(firstNonEmpty(match[1], match[2]))
	}

	return example
}

func firstCodeBlock(content string) (string, string) {
	lines := strings.Split(content, "\n")
	inBlock := false
	var language string
	var body []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !inBlock {
			if strings.HasPrefix(trimmed, "```") {
				inBlock = true
				language = strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
			}
			continue
		}

		if strings.HasPrefix(trimmed, "```") {
			return language, strings.TrimSpace(strings.Join(body, "\n"))
		}

		body = append(body, line)
	}

	return "", ""
}

func buildAbsoluteEndpointURL(baseURL, exampleURL, endpointPath string) (string, error) {
	switch {
	case strings.TrimSpace(baseURL) != "":
		return buildAbsoluteURLFromBase(baseURL, endpointPath)
	case strings.TrimSpace(exampleURL) != "":
		return buildAbsoluteURLFromExample(exampleURL, endpointPath)
	default:
		return "", ErrMarkdownBaseURLNotFound
	}
}

func buildAbsoluteURLFromBase(baseURL, endpointPath string) (string, error) {
	parsedBase, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil || parsedBase.Scheme == "" || parsedBase.Host == "" {
		return "", ErrInvalidMarkdownDocument
	}

	baseSegments := splitPathSegments(parsedBase.Path)
	templateSegments := splitPathSegments(endpointPath)
	if len(baseSegments) > 0 && len(templateSegments) > 0 && baseSegments[len(baseSegments)-1] == templateSegments[0] {
		baseSegments = baseSegments[:len(baseSegments)-1]
	}

	parsedBase.Path = joinPathSegments(append(baseSegments, templateSegments...))
	parsedBase.RawPath = ""
	parsedBase.RawQuery = ""

	return parsedBase.String(), nil
}

func buildAbsoluteURLFromExample(exampleURL, endpointPath string) (string, error) {
	parsedExample, err := url.Parse(strings.TrimSpace(exampleURL))
	if err != nil || parsedExample.Scheme == "" || parsedExample.Host == "" {
		return "", ErrInvalidMarkdownDocument
	}

	exampleSegments := splitPathSegments(parsedExample.Path)
	templateSegments := splitPathSegments(endpointPath)
	start, ok := findTemplateMatch(exampleSegments, templateSegments)
	prefixSegments := make([]string, 0)
	if ok {
		prefixSegments = append(prefixSegments, exampleSegments[:start]...)
	}

	parsedExample.Path = joinPathSegments(append(prefixSegments, templateSegments...))
	parsedExample.RawPath = ""
	parsedExample.RawQuery = ""

	return parsedExample.String(), nil
}

func splitPathSegments(path string) []string {
	trimmed := strings.Trim(strings.TrimSpace(path), "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func joinPathSegments(segments []string) string {
	if len(segments) == 0 {
		return "/"
	}
	return "/" + strings.Join(segments, "/")
}

func findTemplateMatch(exampleSegments, templateSegments []string) (int, bool) {
	if len(templateSegments) == 0 || len(exampleSegments) < len(templateSegments) {
		return 0, false
	}

	for start := 0; start <= len(exampleSegments)-len(templateSegments); start++ {
		matched := true
		for index, templateSegment := range templateSegments {
			if isPathPlaceholder(templateSegment) {
				continue
			}
			if exampleSegments[start+index] != templateSegment {
				matched = false
				break
			}
		}
		if matched {
			return start, true
		}
	}

	return 0, false
}

func extractPathParamsFromExample(endpointPath, exampleURL string) map[string]string {
	if strings.TrimSpace(exampleURL) == "" {
		return nil
	}

	parsedExample, err := url.Parse(strings.TrimSpace(exampleURL))
	if err != nil {
		return nil
	}

	exampleSegments := splitPathSegments(parsedExample.Path)
	templateSegments := splitPathSegments(endpointPath)
	start, ok := findTemplateMatch(exampleSegments, templateSegments)
	if !ok {
		return nil
	}

	result := make(map[string]string)
	for index, templateSegment := range templateSegments {
		if !isPathPlaceholder(templateSegment) {
			continue
		}

		value := exampleSegments[start+index]
		if !isConcretePathValue(templateSegment, value) {
			continue
		}

		result[strings.TrimPrefix(templateSegment, ":")] = value
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func isPathPlaceholder(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

func isConcretePathValue(templateSegment, exampleSegment string) bool {
	switch {
	case exampleSegment == "":
		return false
	case exampleSegment == templateSegment:
		return false
	case strings.HasPrefix(exampleSegment, ":"):
		return false
	case strings.HasPrefix(exampleSegment, "{{") && strings.HasSuffix(exampleSegment, "}}"):
		return false
	default:
		return true
	}
}

func parseQueryParamsFromExample(exampleURL string) []request.KeyValue {
	if strings.TrimSpace(exampleURL) == "" {
		return nil
	}

	parsedExample, err := url.Parse(strings.TrimSpace(exampleURL))
	if err != nil || parsedExample.RawQuery == "" {
		return nil
	}

	params := make([]request.KeyValue, 0, len(parsedExample.Query()))
	for key, values := range parsedExample.Query() {
		for _, value := range values {
			params = append(params, request.KeyValue{
				Key:     key,
				Value:   value,
				Enabled: true,
			})
		}
	}

	if len(params) == 0 {
		return nil
	}

	return params
}

func filterNonAuthHeaders(headers []request.KeyValue) []request.KeyValue {
	filtered := make([]request.KeyValue, 0, len(headers))
	for _, header := range headers {
		if strings.EqualFold(strings.TrimSpace(header.Key), "Authorization") {
			continue
		}
		filtered = append(filtered, header)
	}
	return filtered
}

func deriveBodyType(codeBlockLang string, headers []request.KeyValue) string {
	if strings.EqualFold(strings.TrimSpace(codeBlockLang), "json") {
		return "json"
	}

	for _, header := range headers {
		if !strings.EqualFold(strings.TrimSpace(header.Key), "Content-Type") {
			continue
		}
		if strings.Contains(strings.ToLower(header.Value), "application/json") {
			return "json"
		}
	}

	return "none"
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
