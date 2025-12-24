import os
from dotenv import load_dotenv

# Load environment variables from .env file in the parent directory
load_dotenv(os.path.join(os.path.dirname(__file__), '..', '.env'))

def get_mysql_config(database=None):
    """
    Returns MySQL configuration dictionary from environment variables.
    Optionally overrides the database name.
    """
    config = {
        'host': os.getenv('MYSQL_HOST', 'localhost'),
        'port': int(os.getenv('MYSQL_PORT', 3306)),
        'user': os.getenv('MYSQL_USER', 'root'),
        'password': os.getenv('MYSQL_PASSWORD', ''),
        'charset': 'utf8mb4',
        'use_unicode': True
    }
    
    if database:
        config['database'] = database
    elif os.getenv('MYSQL_DATABASE'):
         config['database'] = os.getenv('MYSQL_DATABASE')
         
    return config

def get_tw_world_db():
    return os.getenv('MYSQL_DATABASE', 'tw_world')

def get_aowow_db():
    return os.getenv('AOWOW_DATABASE', 'aowow')
