#!/bin/bash

echo "🚀 Setting up Python application with SQL migrations..."

# Create virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Copy environment file
cp .env.example .env

# Initialize database migrations
python cli.py db init

# Create initial migration
python cli.py db migrate -m "Initial migration"

# Run migrations
python cli.py db upgrade

echo "✅ Setup complete!"
echo ""
echo "Available commands:"
echo "  python cli.py --help                 # Show all commands"
echo "  python cli.py db migrate -m 'msg'    # Create migration"
echo "  python cli.py db upgrade             # Run migrations"
echo "  python cli.py user create -u name -e email@example.com"
echo "  python cli.py user list              # List users"
echo "  python cli.py status                 # App status"
echo "  python cli.py shell                  # Interactive shell"
