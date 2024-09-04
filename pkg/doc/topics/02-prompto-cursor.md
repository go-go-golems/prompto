---
Title: Setting up Prompto for Cursor Development
Slug: cursor-setup
Short: A guide to setting up and using Prompto with Cursor IDE for efficient development workflows
Topics:
- setup
- cursor
- development
- ngrok
Commands:
- prompto serve
- prompto config repositories add
Flags:
- --all
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

# Setting up Prompto for Cursor Development

This tutorial will guide you through setting up prompto and catter for use with Cursor IDE, allowing you to efficiently manage and use prompts in your development workflow.

## Step 1: Install Catter

Catter is a useful tool for preparing codebase content for LLM contexts. Install it using the following command:

```bash
go install github.com/go-go-golems/go-go-labs/cmd/apps/catter@latest
```

## Step 2: Install Prompto

Prompto is a command-line tool for generating context for prompting Large Language Models. You can install it using one of the following methods:

### Using Homebrew:

```bash
brew tap go-go-golems/go-go-go
brew install go-go-golems/go-go-go/prompto
```

### Using apt-get:

```bash
echo "deb [trusted=yes] https://apt.fury.io/go-go-golems/ /" >> /etc/apt/sources.list.d/fury.list
apt-get update
apt-get install prompto
```

### Using yum:

```bash
echo "
[fury]
name=Gemfury Private Repo
baseurl=https://yum.fury.io/go-go-golems/
enabled=1
gpgcheck=0
" >> /etc/yum.repos.d/fury.repo
yum install prompto
```

### Using go get:

```bash
go get -u github.com/go-go-golems/prompto/cmd/prompto
```

## Step 3: Add Your Repository to Prompto

To add your development repository to prompto, use the following command:

```bash
prompto config repositories add /path/to/your/repository
```

Replace `/path/to/your/repository` with the actual path to your development directory.

## Step 4: Write Dynamic Prompts

Create a `prompto` directory in your repository:

```bash
mkdir /path/to/your/repository/prompto
```

Write your dynamic prompts as executable scripts in this directory. For example:

```bash
vim /path/to/your/repository/prompto/example-prompt.sh
```

Make the script executable:

```bash
chmod +x /path/to/your/repository/prompto/example-prompt.sh
```

Remember, these scripts will be executed in your repository's working directory.

## Step 5: Run Prompto Serve

Start the prompto web server:

```bash
prompto serve
```

This will typically start a local web server on `http://localhost:8080`.

## Step 6: Setup ngrok for Public Access

Cursor IDE doesn't work well with localhost URLs, so we'll use ngrok to make the prompto server publicly accessible.

1. Install ngrok if you haven't already: https://ngrok.com/download

2. Create or edit your ngrok configuration file:

```bash
vim ~/.config/ngrok/ngrok.yml
```

3. Add the following content (replace `YOUR_AUTH_TOKEN` with your actual ngrok auth token):

```yaml
version: "2"
authtoken: YOUR_AUTH_TOKEN
tunnels:
  webapp1:
    proto: http
    addr: 3000
  webapp2:
    proto: http
    addr: 3001
  webapp3:
    proto: http
    addr: 8080
```

4. Start ngrok:

```bash
ngrok start --all
```

This will create public URLs for your local ports, including the prompto server on port 8080.

## Step 7: Use Prompto URLs in Cursor IDE

1. Open the prompto web interface using the public URL provided by ngrok (e.g., `https://abcd1234.ngrok.io`).

2. Navigate to the prompt you want to use and copy its URL.

3. In Cursor IDE's chat or composer, paste the copied URL as additional context.

Now you can easily access your prompto prompts directly from Cursor IDE, enhancing your development workflow with relevant context.