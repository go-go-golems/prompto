# Prompto Web Server Specification

## Small features

- show that the server is running on stdout
- log the hits as concise log line on stdout
- read the list of repositories first and don't reload it on each request.

## Software engineering patterns

- create a struct to keep the server state
- implement htmx backend endpoints rendering html using go template
- create a seaprate file for each API endpoint handler

## Behaviour + API

Prompto is a web server that renders context files and provides a directory-based structure for accessing and searching these contexts. It supports both plain text and HTML responses, with a focus on easy navigation and search functionality.

## Core Functionality

1. **Directory Structure**: 
   - The server organizes content in a directory/name structure (e.g., bubbles/bubble, system/info).
   - Each directory can contain multiple named contexts.

2. **Content Rendering**:
   - The server can render both plain text and HTML content based on the request.
   - For directory listings, it generates an HTML page with search functionality.
   - For specific contexts, it returns the appropriate content type (plain text or HTML).

3. **URL Handling**:
   - Root URL ("/") returns a list of all available top-level directories.
   - Directory URLs (e.g., "/bubbletea/") return a list of contents within that directory.
   - Specific context URLs (e.g., "/bubbles/bubble") return the content of that context.

4. **Search Functionality**:
   - Implements a search feature for directory listings using HTMX.
   - Provides server-side filtering of directory contents based on search queries.

## Detailed Specifications

### 1. URL Handling and Responses

#### Root URL ("/")
- **Request**: GET /
- **Response**: 
  - Content-Type: text/html
  - Body: HTML page listing all top-level directories with search functionality

#### Directory URL (e.g., "/bubbletea/")
- **Request**: GET /directory-name/
- **Response**:
  - Content-Type: text/html
  - Body: HTML page listing contents of the specified directory with search functionality

#### Specific Context URL (e.g., "/bubbles/bubble")
- **Request**: GET /directory-name/context-name
- **Response**:
  - Content-Type: Depends on the context (text/plain or text/html)
  - Body: Content of the specified context

### 2. HTML Structure for Directory Listings

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>[Directory Name] Contents</title>
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
    <style>
        /* Include the CSS styles here */
    </style>
</head>
<body>
    <h1>[Directory Name] Contents</h1>
    
    <input type="text" 
           placeholder="Search contents..." 
           name="search" 
           hx-post="/[directory-name]/search"
           hx-trigger="keyup changed delay:200ms, search"
           hx-target="#content-list"
           hx-indicator=".htmx-indicator">
    
    <div class="htmx-indicator">Searching...</div>
    
    <ul id="content-list">
        <!-- List items here -->
    </ul>
</body>
</html>
```

### 3. Search Functionality

#### Client-side (HTMX):
- Triggers a POST request to "/[directory-name]/search" on keyup events in the search input.
- 200ms delay to prevent excessive requests.
- Updates the #content-list with the server response.

#### Server-side:
- **Endpoint**: POST /[directory-name]/search
- **Request Body**: search query (form-encoded)
- **Response**: 
  - Content-Type: text/html
  - Body: Filtered list of `<li>` elements matching the search query

## Additional API 

- /refresh endpoint to reload the promptos
- /repositories to list the repositories


## New features

- Add a clipboard to the resulting bubble