#!/bin/bash
# Add MIT license header to Go files that don't have one

LICENSE_HEADER='// Copyright (c) 2026 Clotho contributors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

'

find . -name "*.go" -not -path "./vendor/*" | while read -r file; do
    if ! grep -q "^// Copyright" "$file"; then
        # Get line number of package declaration
        package_line=$(grep -n "^package " "$file" | head -1 | cut -d: -f1)
        if [ -n "$package_line" ]; then
            # Get the content before package line
            head -n $((package_line - 1)) "$file" > "$file.tmp"
            # Add license
            echo "$LICENSE_HEADER" >> "$file.tmp"
            # Add the rest of the file starting from package line
            tail -n +$package_line "$file" >> "$file.tmp"
            mv "$file.tmp" "$file"
            echo "Added license to: $file"
        fi
    fi
done
