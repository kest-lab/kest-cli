import { MetadataRoute } from 'next';

const BASE_URL = 'http://localhost:3000';

export default function robots(): MetadataRoute.Robots {
  return {
    rules: [{ userAgent: '*', allow: '/' }],
    sitemap: `${BASE_URL}/sitemap.xml`,
  };
}
