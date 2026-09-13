import fs from "node:fs";
import path from "node:path";
import { marked } from "marked";

const service = process.argv[2];

if (!service) {
  throw new Error("Service ID is required.");
}

const sourceDocsDir = `src/api/services/${service}/docs`;
const outputDocsDir = `dist/api/oas/${service}/docs`;
const templatePath = "tools/templates/external-doc.html";

if (!fs.existsSync(sourceDocsDir)) {
  console.log(`No externalDocs: ${sourceDocsDir}`);
  process.exit(0);
}

const template = fs.readFileSync(templatePath, "utf8");

const markdownFiles = fs
  .readdirSync(sourceDocsDir)
  .filter((file) => path.extname(file).toLowerCase() === ".md");

if (markdownFiles.length === 0) {
  console.log(`No externalDocs: ${sourceDocsDir}`);
  process.exit(0);
}

fs.mkdirSync(outputDocsDir, { recursive: true });

for (const file of markdownFiles) {
  const markdownPath = path.join(sourceDocsDir, file);
  const markdown = fs.readFileSync(markdownPath, "utf8");

  const titleMatch = markdown.match(/^#\s+(.+)$/m);

  if (!titleMatch) {
    throw new Error(`H1 title not found: ${markdownPath}`);
  }

  const title = titleMatch[1].trim();
  const content = marked.parse(markdown);

  const html = template
    .replace("{{title}}", title)
    .replace("{{content}}", content);

  const outputFile = `${path.basename(file, ".md")}.html`;
  const outputPath = path.join(outputDocsDir, outputFile);

  fs.writeFileSync(outputPath, html, "utf8");

  console.log(`Generated: ${outputPath}`);
}
