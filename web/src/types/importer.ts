export interface ImportPostmanCollectionRequest {
  file: File;
  parent_id?: number;
}

export interface ImportPostmanCollectionResponse {
  message?: string;
}

export interface ImportMarkdownCollectionRequest {
  file: File;
  parent_id?: number;
}

export interface ImportMarkdownCollectionModuleResult {
  name: string;
  collection_id: number;
  request_count: number;
}

export interface ImportMarkdownCollectionResponse {
  root_folder_id: number;
  root_folder_name: string;
  collections_created: number;
  requests_created: number;
  modules: ImportMarkdownCollectionModuleResult[];
}
