import click
import os
import subprocess
from sqlalchemy.orm import Session
from database import get_db, User, Post, create_tables

@click.group()
def cli():
    """My App CLI - Database and application management"""
    pass

# Database Management Commands
@cli.group()
def db():
    """Database management commands"""
    pass

@db.command()
def init():
    """Initialize the database with Alembic"""
    click.echo("Initializing database migrations...")
    if not os.path.exists('migrations'):
        subprocess.run(['alembic', 'init', 'migrations'], check=True)
        click.echo("✅ Database migrations initialized!")
    else:
        click.echo("⚠️  Migrations directory already exists")

@db.command()
@click.option('--message', '-m', required=True, help='Migration message')
def migrate(message):
    """Create a new migration"""
    click.echo(f"Creating migration: {message}")
    subprocess.run(['alembic', 'revision', '--autogenerate', '-m', message], check=True)
    click.echo("✅ Migration created!")

@db.command()
def upgrade():
    """Run database migrations"""
    click.echo("Running database migrations...")
    subprocess.run(['alembic', 'upgrade', 'head'], check=True)
    click.echo("✅ Database upgraded!")

@db.command()
def downgrade():
    """Downgrade database by one migration"""
    click.echo("Downgrading database...")
    subprocess.run(['alembic', 'downgrade', '-1'], check=True)
    click.echo("✅ Database downgraded!")

@db.command()
def reset():
    """Reset the database (WARNING: This will delete all data!)"""
    if click.confirm('This will delete all data. Are you sure?'):
        # Remove database file if using SQLite
        db_file = 'app.db'
        if os.path.exists(db_file):
            os.remove(db_file)
        
        # Remove migrations
        if os.path.exists('migrations'):
            import shutil
            shutil.rmtree('migrations')
        
        click.echo("✅ Database reset!")

@db.command()
def create():
    """Create database tables (for development only)"""
    create_tables()
    click.echo("✅ Tables created!")

@db.command()
def status():
    """Show migration status"""
    subprocess.run(['alembic', 'current'], check=True)

# User Management Commands
@cli.group()
def user():
    """User management commands"""
    pass

@user.command()
@click.option('--username', '-u', required=True, help='Username')
@click.option('--email', '-e', required=True, help='Email address')
def create(username, email):
    """Create a new user"""
    db = next(get_db())
    try:
        new_user = User(username=username, email=email)
        db.add(new_user)
        db.commit()
        click.echo(f"✅ User '{username}' created with ID: {new_user.id}")
    except Exception as e:
        db.rollback()
        click.echo(f"❌ Error creating user: {e}")
    finally:
        db.close()

@user.command()
def list():
    """List all users"""
    db = next(get_db())
    try:
        users = db.query(User).all()
        if users:
            click.echo("\n📋 Users:")
            for user in users:
                click.echo(f"  ID: {user.id} | Username: {user.username} | Email: {user.email}")
        else:
            click.echo("No users found")
    finally:
        db.close()

@user.command()
@click.argument('user_id', type=int)
def delete(user_id):
    """Delete a user by ID"""
    db = next(get_db())
    try:
        user = db.query(User).filter(User.id == user_id).first()
        if user:
            username = user.username
            db.delete(user)
            db.commit()
            click.echo(f"✅ User '{username}' deleted")
        else:
            click.echo(f"❌ User with ID {user_id} not found")
    finally:
        db.close()

# Application Commands
@cli.command()
def serve():
    """Start the application server"""
    click.echo("🚀 Starting application server...")
    # Add your server startup code here
    # For example: app.run() if using Flask
    click.echo("Server would start here (implement your server logic)")

@cli.command()
@click.option('--env', default='development', help='Environment to run in')
def status(env):
    """Show application status"""
    click.echo(f"📊 Application Status")
    click.echo(f"Environment: {env}")
    click.echo(f"Database URL: {os.getenv('DATABASE_URL', 'sqlite:///app.db')}")
    
    # Check database connection
    try:
        db = next(get_db())
        user_count = db.query(User).count()
        post_count = db.query(Post).count()
        click.echo(f"Database: ✅ Connected")
        click.echo(f"Users: {user_count}")
        click.echo(f"Posts: {post_count}")
        db.close()
    except Exception as e:
        click.echo(f"Database: ❌ Error - {e}")

@cli.command()
def shell():
    """Start an interactive Python shell with app context"""
    import code
    from database import engine, SessionLocal, User, Post
    
    banner = """
🐍 Interactive Shell
Available objects:
  - engine: SQLAlchemy engine
  - SessionLocal: Session factory
  - User, Post: Models
  - db = SessionLocal() to get a database session
"""
    
    local_vars = {
        'engine': engine,
        'SessionLocal': SessionLocal,
        'User': User,
        'Post': Post,
        'db': SessionLocal()
    }
    
    code.interact(banner=banner, local=local_vars)

if __name__ == '__main__':
    cli()
