import fs from 'fs';

const versionFilePath = 'cmd/version.go';

// Read the current content of the version.go file
const content = fs.readFileSync(versionFilePath, 'utf8');

// Regular expression to find the version constant
const versionRegex = /const version = "([^"]+)"/;

// Get the new version from command line arguments
const newVersion = process.argv[2];

if (!newVersion) {
    console.error('Usage: node update-version.mjs <new-version>');
    process.exit(1);
}

// Replace the version
const updatedContent = content.replace(versionRegex, `const version = "${newVersion}"`);

// Write the updated content back to the file
fs.writeFileSync(versionFilePath, updatedContent);

console.log(`Version updated to: ${newVersion}`);