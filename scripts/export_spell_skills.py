import db_config
"""
Export spell skills from aowow database for spell categorization
Creates:
- spell_skill_categories: Category groupings (Class Skills, Professions, etc.)
- spell_skills: Individual skills (Frost, Fire, Alchemy, etc.)  
- spell_skill_spells: Mapping of spells to skills
"""
import mysql.connector
import sqlite3
import os

# MySQL config
mysql_config = db_config.get_mysql_config()

# SQLite path
SQLITE_DB = os.path.join(os.path.dirname(__file__), '..', 'data', 'shelllab.db')

def main():
    print("Connecting to MySQL (aowow)...")
    mysql_conn = mysql.connector.connect(**mysql_config)
    mysql_cursor = mysql_conn.cursor(dictionary=True)
    
    print("Connecting to SQLite...")
    sqlite_conn = sqlite3.connect(SQLITE_DB)
    sqlite_cursor = sqlite_conn.cursor()
    
    # Create tables
    print("Creating spell skill tables...")
    sqlite_cursor.executescript("""
        DROP TABLE IF EXISTS spell_skill_categories;
        DROP TABLE IF EXISTS spell_skills;
        DROP TABLE IF EXISTS spell_skill_spells;
        
        CREATE TABLE spell_skill_categories (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL
        );
        
        CREATE TABLE spell_skills (
            id INTEGER PRIMARY KEY,
            category_id INTEGER NOT NULL,
            name TEXT NOT NULL,
            FOREIGN KEY (category_id) REFERENCES spell_skill_categories(id)
        );
        
        CREATE TABLE spell_skill_spells (
            skill_id INTEGER NOT NULL,
            spell_id INTEGER NOT NULL,
            PRIMARY KEY (skill_id, spell_id),
            FOREIGN KEY (skill_id) REFERENCES spell_skills(id)
        );
        
        CREATE INDEX idx_spell_skill_spells_spell ON spell_skill_spells(spell_id);
    """)
    
    # Category names mapping
    category_names = {
        6: "Weapon Skills",
        7: "Class Skills", 
        8: "Armor Skills",
        9: "Secondary Skills",
        10: "Languages",
        11: "Professions",
        12: "Generic"
    }
    
    # Insert categories
    print("Inserting spell skill categories...")
    for cat_id, cat_name in category_names.items():
        sqlite_cursor.execute(
            "INSERT INTO spell_skill_categories (id, name) VALUES (?, ?)",
            (cat_id, cat_name)
        )
    
    # Get skills from aowow
    print("Fetching skills from aowow...")
    mysql_cursor.execute("""
        SELECT skillID, categoryID, name_loc0 
        FROM aowow_skill 
        WHERE categoryID >= 6 AND categoryID <= 12
        ORDER BY categoryID, name_loc0
    """)
    skills = mysql_cursor.fetchall()
    
    print(f"Inserting {len(skills)} skills...")
    for skill in skills:
        sqlite_cursor.execute(
            "INSERT OR IGNORE INTO spell_skills (id, category_id, name) VALUES (?, ?, ?)",
            (skill['skillID'], skill['categoryID'], skill['name_loc0'])
        )
    
    # Get skill-spell mappings from aowow
    print("Fetching skill-spell mappings...")
    mysql_cursor.execute("""
        SELECT skillID, spellID 
        FROM aowow_skill_line_ability
    """)
    mappings = mysql_cursor.fetchall()
    
    print(f"Inserting {len(mappings)} spell-skill mappings...")
    for mapping in mappings:
        sqlite_cursor.execute(
            "INSERT OR IGNORE INTO spell_skill_spells (skill_id, spell_id) VALUES (?, ?)",
            (mapping['skillID'], mapping['spellID'])
        )
    
    sqlite_conn.commit()
    
    # Verify
    sqlite_cursor.execute("SELECT COUNT(*) FROM spell_skill_categories")
    cat_count = sqlite_cursor.fetchone()[0]
    sqlite_cursor.execute("SELECT COUNT(*) FROM spell_skills")
    skill_count = sqlite_cursor.fetchone()[0]
    sqlite_cursor.execute("SELECT COUNT(*) FROM spell_skill_spells")
    mapping_count = sqlite_cursor.fetchone()[0]
    
    print(f"\nDone! Imported:")
    print(f"  - {cat_count} categories")
    print(f"  - {skill_count} skills")
    print(f"  - {mapping_count} spell-skill mappings")
    
    mysql_cursor.close()
    mysql_conn.close()
    sqlite_cursor.close()
    sqlite_conn.close()

if __name__ == "__main__":
    main()
