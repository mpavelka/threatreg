#!/usr/bin/env python3
"""
Service Documentation Generator

This script scans all Go service files in the internal/service directory,
extracts public function documentation, and generates a comprehensive HTML
documentation with expandable sections organized by service category.
"""

import os
import re
import json
from pathlib import Path
from datetime import datetime
from typing import Dict, List, Tuple, Optional


def extract_go_functions(file_path: str) -> List[Dict]:
    """
    Extract public functions and their documentation from Go files.
    
    Args:
        file_path: Path to the Go file to analyze
        
    Returns:
        List of dictionaries containing function information
    """
    functions = []
    
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        print(f"Error reading {file_path}: {e}")
        return functions
    
    # Regular expression to match function with optional documentation
    # Match comments followed by function signature
    func_pattern = r'(?:^//.*\n)*^func\s+([A-Z][a-zA-Z0-9]*)\s*\([^)]*\)(?:\s*\([^)]*\))?\s*(?:\*?[a-zA-Z0-9\[\]\.]+)?\s*(?:,\s*error)?\s*{'
    
    lines = content.split('\n')
    i = 0
    
    while i < len(lines):
        line = lines[i].strip()
        
        # Look for function definitions starting with capital letter
        if line.startswith('func ') and re.match(r'func\s+[A-Z]', line):
            # Extract documentation comments above the function
            docs = []
            j = i - 1
            
            # Go backwards to collect documentation comments
            while j >= 0 and (lines[j].strip().startswith('//') or lines[j].strip() == ''):
                if lines[j].strip().startswith('//'):
                    # Remove the // prefix and clean up
                    doc_line = lines[j].strip()[2:].strip()
                    docs.insert(0, doc_line)
                elif lines[j].strip() == '':
                    # Empty line - continue looking for more docs
                    pass
                else:
                    break
                j -= 1
            
            # Extract function signature (may span multiple lines)
            signature_lines = []
            k = i
            brace_count = 0
            paren_count = 0
            in_func = False
            
            while k < len(lines):
                curr_line = lines[k].strip()
                if not in_func and curr_line.startswith('func '):
                    in_func = True
                
                if in_func:
                    signature_lines.append(lines[k].rstrip())
                    
                    # Count parentheses and braces to find end of signature
                    paren_count += curr_line.count('(') - curr_line.count(')')
                    
                    if '{' in curr_line:
                        brace_count += curr_line.count('{')
                        break
                    
                    if paren_count == 0 and k > i:
                        # Look for opening brace on next lines
                        next_k = k + 1
                        while next_k < len(lines) and lines[next_k].strip() == '':
                            next_k += 1
                        if next_k < len(lines) and lines[next_k].strip().startswith('{'):
                            break
                
                k += 1
            
            # Clean up and extract function name and signature
            full_signature = ' '.join(signature_lines).strip()
            
            # Extract function name
            name_match = re.match(r'func\s+([A-Z][a-zA-Z0-9]*)', full_signature)
            if name_match:
                func_name = name_match.group(1)
                
                # Clean up signature
                clean_signature = re.sub(r'\s+', ' ', full_signature)
                clean_signature = re.sub(r'\s*{\s*$', '', clean_signature)
                
                functions.append({
                    'name': func_name,
                    'signature': clean_signature,
                    'documentation': docs,
                    'line_number': i + 1
                })
        
        i += 1
    
    return functions


def categorize_service_file(file_path: str) -> str:
    """
    Determine the service category based on file path and name.
    
    Args:
        file_path: Path to the service file
        
    Returns:
        Category name for the service
    """
    path_obj = Path(file_path)
    
    # Handle threat_pattern subdirectory
    if 'threat_pattern' in path_obj.parts:
        filename = path_obj.stem
        if 'pattern_condition' in filename:
            return 'Threat Pattern Conditions'
        elif 'instance_threat_pattern' in filename:
            return 'Instance Threat Pattern Evaluation'
        else:
            return 'Threat Pattern Management'
    
    # Main service files
    filename = path_obj.stem
    category_map = {
        'control': 'Control Management',
        'domain': 'Domain Management', 
        'instance': 'Instance Management',
        'product': 'Product Management',
        'threat': 'Threat Management',
        'relationship': 'Relationship Management',
        'tag': 'Tag Management',
        'threat_resolution': 'Threat Resolution Management'
    }
    
    return category_map.get(filename, 'Miscellaneous')


def scan_service_directory(service_dir: str) -> Dict[str, List[Dict]]:
    """
    Scan the service directory and extract all public functions organized by category.
    
    Args:
        service_dir: Path to the internal/service directory
        
    Returns:
        Dictionary mapping category names to lists of functions
    """
    categories = {}
    
    # Find all Go files (excluding test files)
    for root, dirs, files in os.walk(service_dir):
        for file in files:
            if file.endswith('.go') and not file.endswith('_test.go'):
                file_path = os.path.join(root, file)
                category = categorize_service_file(file_path)
                
                functions = extract_go_functions(file_path)
                
                if functions:
                    if category not in categories:
                        categories[category] = {}
                    
                    # Use relative path as key for file organization
                    rel_path = os.path.relpath(file_path, service_dir)
                    categories[category][rel_path] = functions
    
    return categories


def generate_html_documentation(categories: Dict[str, Dict[str, List[Dict]]], output_file: str):
    """
    Generate HTML documentation with expandable sections.
    
    Args:
        categories: Dictionary of service categories and their functions
        output_file: Path to output HTML file
    """
    
    total_functions = sum(len(funcs) for cat_files in categories.values() 
                         for funcs in cat_files.values())
    
    html_content = f"""<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Threatreg Service Layer Documentation</title>
    <style>
        body {{
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f8f9fa;
        }}
        
        .header {{
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 2rem;
            border-radius: 10px;
            margin-bottom: 2rem;
            text-align: center;
        }}
        
        .header h1 {{
            margin: 0;
            font-size: 2.5rem;
            font-weight: 300;
        }}
        
        .header p {{
            margin: 1rem 0 0 0;
            opacity: 0.9;
            font-size: 1.1rem;
        }}
        
        .stats {{
            background: white;
            padding: 1.5rem;
            border-radius: 8px;
            margin-bottom: 2rem;
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1rem;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }}
        
        .stat-item {{
            text-align: center;
            padding: 1rem;
            border-radius: 6px;
            background: #f8f9fa;
        }}
        
        .stat-number {{
            font-size: 2rem;
            font-weight: bold;
            color: #667eea;
            display: block;
        }}
        
        .category {{
            background: white;
            margin-bottom: 1.5rem;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            overflow: hidden;
        }}
        
        .category-header {{
            background: #667eea;
            color: white;
            padding: 1rem 1.5rem;
            cursor: pointer;
            user-select: none;
            transition: background-color 0.3s ease;
            position: relative;
        }}
        
        .category-header:hover {{
            background: #5a67d8;
        }}
        
        .category-header h2 {{
            margin: 0;
            font-size: 1.3rem;
            font-weight: 500;
        }}
        
        .category-header::after {{
            content: '▼';
            position: absolute;
            right: 1.5rem;
            top: 50%;
            transform: translateY(-50%);
            transition: transform 0.3s ease;
        }}
        
        .category.collapsed .category-header::after {{
            transform: translateY(-50%) rotate(-90deg);
        }}
        
        .category-content {{
            max-height: 1000px;
            overflow: hidden;
            transition: max-height 0.3s ease;
        }}
        
        .category.collapsed .category-content {{
            max-height: 0;
        }}
        
        .file-section {{
            border-bottom: 1px solid #e2e8f0;
        }}
        
        .file-section:last-child {{
            border-bottom: none;
        }}
        
        .file-header {{
            background: #f8f9fa;
            padding: 0.75rem 1.5rem;
            font-weight: 500;
            color: #4a5568;
            font-family: 'Monaco', 'Consolas', monospace;
            font-size: 0.9rem;
        }}
        
        .function {{
            padding: 1.5rem;
            border-bottom: 1px solid #f1f5f9;
        }}
        
        .function:last-child {{
            border-bottom: none;
        }}
        
        .function-name {{
            font-size: 1.2rem;
            font-weight: 600;
            color: #2d3748;
            margin-bottom: 0.5rem;
        }}
        
        .function-signature {{
            background: #f7fafc;
            padding: 1rem;
            border-radius: 6px;
            border-left: 4px solid #667eea;
            font-family: 'Monaco', 'Consolas', monospace;
            font-size: 0.9rem;
            color: #4a5568;
            margin-bottom: 1rem;
            overflow-x: auto;
        }}
        
        .function-docs {{
            color: #4a5568;
            line-height: 1.7;
        }}
        
        .function-docs p {{
            margin: 0 0 0.5rem 0;
        }}
        
        .function-docs p:last-child {{
            margin-bottom: 0;
        }}
        
        .no-docs {{
            color: #a0aec0;
            font-style: italic;
        }}
        
        .footer {{
            text-align: center;
            padding: 2rem;
            color: #718096;
            border-top: 1px solid #e2e8f0;
            margin-top: 3rem;
        }}
        
        .toc {{
            background: white;
            padding: 1.5rem;
            border-radius: 8px;
            margin-bottom: 2rem;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }}
        
        .toc h3 {{
            margin-top: 0;
            color: #2d3748;
        }}
        
        .toc ul {{
            list-style: none;
            padding: 0;
        }}
        
        .toc li {{
            padding: 0.25rem 0;
        }}
        
        .toc a {{
            color: #667eea;
            text-decoration: none;
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            transition: background-color 0.2s ease;
        }}
        
        .toc a:hover {{
            background-color: #f7fafc;
        }}
        
        @media (max-width: 768px) {{
            body {{
                padding: 10px;
            }}
            
            .header {{
                padding: 1.5rem;
            }}
            
            .header h1 {{
                font-size: 2rem;
            }}
            
            .stats {{
                grid-template-columns: 1fr;
            }}
        }}
    </style>
</head>
<body>
    <div class="header">
        <h1>Threatreg Service Layer Documentation</h1>
        <p>Comprehensive documentation for all public service functions</p>
    </div>
    
    <div class="stats">
        <div class="stat-item">
            <span class="stat-number">{len(categories)}</span>
            Service Categories
        </div>
        <div class="stat-item">
            <span class="stat-number">{total_functions}</span>
            Public Functions
        </div>
        <div class="stat-item">
            <span class="stat-number">{datetime.now().strftime('%Y-%m-%d')}</span>
            Generated
        </div>
    </div>
    
    <div class="toc">
        <h3>Table of Contents</h3>
        <ul>
"""
    
    # Generate table of contents
    for category in sorted(categories.keys()):
        func_count = sum(len(funcs) for funcs in categories[category].values())
        html_content += f'            <li><a href="#{category.lower().replace(" ", "-")}">{category}</a> ({func_count} functions)</li>\n'
    
    html_content += """        </ul>
    </div>
"""
    
    # Generate content for each category
    for category in sorted(categories.keys()):
        category_id = category.lower().replace(' ', '-')
        html_content += f"""    
    <div class="category" id="{category_id}">
        <div class="category-header" onclick="toggleCategory(this)">
            <h2>{category}</h2>
        </div>
        <div class="category-content">
"""
        
        # Group functions by file
        for file_path in sorted(categories[category].keys()):
            functions = categories[category][file_path]
            
            html_content += f"""            <div class="file-section">
                <div class="file-header">{file_path}</div>
"""
            
            for func in functions:
                docs_html = ""
                if func['documentation']:
                    docs_text = '\\n'.join(func['documentation'])
                    # Convert to HTML paragraphs
                    paragraphs = docs_text.split('\\n\\n')
                    docs_html = '\\n'.join(f"<p>{para.strip()}</p>" for para in paragraphs if para.strip())
                else:
                    docs_html = '<p class="no-docs">No documentation available.</p>'
                
                html_content += f"""                <div class="function">
                    <div class="function-name">{func['name']}</div>
                    <div class="function-signature">{func['signature']}</div>
                    <div class="function-docs">
                        {docs_html}
                    </div>
                </div>
"""
            
            html_content += "            </div>\\n"
        
        html_content += """        </div>
    </div>
"""
    
    html_content += f"""
    <div class="footer">
        <p>Generated on {datetime.now().strftime('%Y-%m-%d at %H:%M:%S')} | 
           Total functions documented: {total_functions}</p>
    </div>

    <script>
        function toggleCategory(header) {{
            const category = header.parentElement;
            category.classList.toggle('collapsed');
        }}
        
        // Initialize all categories as expanded
        document.addEventListener('DOMContentLoaded', function() {{
            // You can uncomment the next line to start with all categories collapsed
            // document.querySelectorAll('.category').forEach(cat => cat.classList.add('collapsed'));
        }});
    </script>
</body>
</html>"""
    
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(html_content)


def main():
    """Main function to generate service documentation."""
    # Determine paths
    script_dir = Path(__file__).parent
    project_root = script_dir.parent
    service_dir = project_root / 'internal' / 'service'
    output_file = script_dir / 'service_documentation.html'
    
    print("🔍 Scanning service directory for public functions...")
    print(f"Service directory: {service_dir}")
    
    if not service_dir.exists():
        print(f"❌ Service directory not found: {service_dir}")
        return
    
    # Scan and extract functions
    categories = scan_service_directory(str(service_dir))
    
    if not categories:
        print("❌ No public functions found in service directory")
        return
    
    # Print summary
    total_functions = sum(len(funcs) for cat_files in categories.values() 
                         for funcs in cat_files.values())
    
    print(f"\\n📊 Summary:")
    print(f"   • Categories found: {len(categories)}")
    print(f"   • Total functions: {total_functions}")
    
    for category in sorted(categories.keys()):
        func_count = sum(len(funcs) for funcs in categories[category].values())
        print(f"   • {category}: {func_count} functions")
    
    # Generate HTML documentation
    print(f"\\n📝 Generating HTML documentation...")
    generate_html_documentation(categories, str(output_file))
    
    print(f"✅ Documentation generated successfully!")
    print(f"📄 Output file: {output_file}")
    print(f"🌐 Open in browser: file://{output_file.absolute()}")


if __name__ == '__main__':
    main()