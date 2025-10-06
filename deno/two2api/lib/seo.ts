// SEO utilities for OpenAI-compatible API proxy

import type { ProxyConfig } from "./types.ts";

/**
 * Generate SEO meta tags for HTML head
 */
export function getSeoMeta(
  config: ProxyConfig,
  pageTitle: string,
  pageDescription?: string,
  currentUrl?: string
): string {
  const title = config.seoTitle
    ? `${pageTitle} - ${config.seoTitle}`
    : `${pageTitle} - ${config.serviceName}`;
  const description =
    pageDescription || config.seoDescription || `${config.serviceName} - OpenAI兼容的API代理服务`;
  const keywords = config.seoKeywords || `OpenAI,API,Proxy,${config.serviceName},AI,GPT`;
  const author = config.seoAuthor || config.serviceName;
  const ogImage = config.seoOgImage || "";

  return `
    <!-- SEO Meta Tags -->
    <meta name="description" content="${description}">
    <meta name="keywords" content="${keywords}">
    <meta name="author" content="${author}">

    <!-- Open Graph Meta Tags -->
    <meta property="og:type" content="website">
    <meta property="og:title" content="${title}">
    <meta property="og:description" content="${description}">
    ${currentUrl ? `<meta property="og:url" content="${currentUrl}">` : ""}
    ${ogImage ? `<meta property="og:image" content="${ogImage}">` : ""}
    <meta property="og:site_name" content="${config.serviceName}">

    <!-- Twitter Card Meta Tags -->
    <meta name="twitter:card" content="summary_large_image">
    <meta name="twitter:title" content="${title}">
    <meta name="twitter:description" content="${description}">
    ${ogImage ? `<meta name="twitter:image" content="${ogImage}">` : ""}

    <!-- Additional SEO -->
    <meta name="robots" content="index, follow">
    <meta name="googlebot" content="index, follow">
    <link rel="canonical" href="${currentUrl || ""}">
  `;
}

/**
 * Generate structured data (JSON-LD) for SEO
 */
export function getStructuredData(config: ProxyConfig, currentUrl?: string): string {
  const structuredData = {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    name: config.serviceName,
    description: config.seoDescription || `${config.serviceName} - OpenAI兼容的API代理服务`,
    url: currentUrl || config.githubRepo,
    applicationCategory: "DeveloperApplication",
    operatingSystem: "Cross-platform",
    offers: {
      "@type": "Offer",
      price: "0",
      priceCurrency: "USD",
    },
    author: {
      "@type": "Organization",
      name: config.seoAuthor || config.serviceName,
      url: config.githubRepo,
    },
  };

  return `
    <script type="application/ld+json">
    ${JSON.stringify(structuredData, null, 2)}
    </script>
  `;
}

/**
 * Get default SEO config values
 */
export function getDefaultSeoConfig() {
  return {
    seoTitle: "",
    seoDescription: "OpenAI兼容的API代理服务，支持流式和非流式响应",
    seoKeywords: "OpenAI,API,Proxy,AI,GPT,Deno,TypeScript",
    seoAuthor: "OpenAI API Proxy",
    seoOgImage: "",
  };
}
