const { execSync } = require('child_process');
const { existsSync, mkdirSync, createWriteStream, chmodSync } = require('fs');
const { get } = require('https');
const { join } = require('path');

const BIN_DIR = join(__dirname, 'bin');
const BIN_NAME = process.platform === 'win32' ? 'autocommit-cli.exe' : 'autocommit-cli';
const BIN_PATH = join(BIN_DIR, BIN_NAME);

const REPO_URL = 'https://github.com/urstruelysv/autocommit-cli';
const VERSION = 'v0.1.0'; // This should ideally be read from package.json

function getDownloadUrl() {
  const platform = process.platform;
  const arch = process.arch;

  let os = '';
  let architecture = '';

  if (platform === 'darwin') {
    os = 'darwin';
  } else if (platform === 'linux') {
    os = 'linux';
  } else if (platform === 'win32') {
    os = 'windows';
  } else {
    console.error(`Unsupported OS: ${platform}`);
    process.exit(1);
  }

  if (arch === 'x64') {
    architecture = 'amd64';
  } else if (arch === 'arm64') {
    architecture = 'arm64';
  } else {
    console.error(`Unsupported architecture: ${arch}`);
    process.exit(1);
  }

  // Example: autocommit-cli-darwin-amd64.tar.gz or autocommit-cli-windows-amd64.zip
  const ext = os === 'windows' ? 'zip' : 'tar.gz';
  return `${REPO_URL}/releases/download/${VERSION}/autocommit-cli-${os}-${architecture}.${ext}`;
}

function downloadBinary(url, dest) {
  console.log(`Downloading from: ${url}`);
  if (!existsSync(BIN_DIR)) {
    mkdirSync(BIN_DIR, { recursive: true });
  }

  const file = createWriteStream(dest);
  get(url, (response) => {
    if (response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
      // Handle redirects
      console.log(`Redirecting to: ${response.headers.location}`);
      downloadBinary(response.headers.location, dest);
      return;
    }
    if (response.statusCode !== 200) {
      console.error(`Failed to download binary. Status Code: ${response.statusCode}`);
      process.exit(1);
    }
    response.pipe(file);
    file.on('finish', () => {
      file.close(() => {
        console.log('Download complete.');
        // Make the binary executable on non-Windows systems
        if (process.platform !== 'win32') {
          chmodSync(dest, '755');
        }
        console.log(`autocommit-cli installed to: ${dest}`);
      });
    });
  }).on('error', (err) => {
    console.error(`Error during download: ${err.message}`);
    process.exit(1);
  });
}

function extractBinary(archivePath, destDir) {
  console.log(`Extracting ${archivePath} to ${destDir}`);
  try {
    if (archivePath.endsWith('.zip')) {
      // For zip files, we need an unzip utility
      execSync(`unzip -o ${archivePath} -d ${destDir}`);
    } else if (archivePath.endsWith('.tar.gz')) {
      execSync(`tar -xzf ${archivePath} -C ${destDir}`);
    } else {
      console.error(`Unsupported archive type: ${archivePath}`);
      process.exit(1);
    }
    console.log('Extraction complete.');
  } catch (error) {
    console.error(`Error during extraction: ${error.message}`);
    process.exit(1);
  }
}

async function main() {
  const downloadUrl = getDownloadUrl();
  const archivePath = join(BIN_DIR, `autocommit-cli-archive.${downloadUrl.split('.').pop()}`);

  await new Promise((resolve) => {
    const file = createWriteStream(archivePath);
    get(downloadUrl, (response) => {
      if (response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
        // Handle redirects
        console.log(`Redirecting to: ${response.headers.location}`);
        get(response.headers.location, (redirectResponse) => {
          redirectResponse.pipe(file);
          file.on('finish', () => file.close(resolve));
        }).on('error', (err) => {
          console.error(`Error during redirect download: ${err.message}`);
          process.exit(1);
        });
        return;
      }
      if (response.statusCode !== 200) {
        console.error(`Failed to download binary. Status Code: ${response.statusCode}`);
        process.exit(1);
      }
      response.pipe(file);
      file.on('finish', () => file.close(resolve));
    }).on('error', (err) => {
      console.error(`Error during download: ${err.message}`);
      process.exit(1);
    });
  });

  extractBinary(archivePath, BIN_DIR);

  // Assuming the extracted binary is directly named 'autocommit-cli' or 'autocommit-cli.exe'
  // and is placed in BIN_DIR
  if (process.platform !== 'win32') {
    chmodSync(BIN_PATH, '755');
  }
  console.log(`autocommit-cli installed to: ${BIN_PATH}`);
}

main();
