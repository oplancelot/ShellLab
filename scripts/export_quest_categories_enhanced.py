import db_config
"""
Export enhanced quest categories from tw_world database
Uses ZoneOrSort to categorize quests:
- Positive values: Zone-based (Barrens, Stranglethorn, etc.)
- Negative values: Class/Profession quests (Warlock, Cooking, etc.)

Creates a comprehensive quest_categories table with proper grouping
"""
import mysql.connector
import sqlite3
import os

# MySQL config
mysql_config = db_config.get_mysql_config()

# SQLite path
SQLITE_DB = os.path.join(os.path.dirname(__file__), '..', 'data', 'shelllab.db')

# Quest category mapping for negative ZoneOrSort values
QUEST_SPECIAL_CATEGORIES = {
    # Class quests (grouped under "Class Quests")
    -61: ("Class Quests", "Warlock"),
    -81: ("Class Quests", "Warrior"),
    -82: ("Class Quests", "Shaman"),
    -141: ("Class Quests", "Paladin"),
    -161: ("Class Quests", "Mage"),
    -162: ("Class Quests", "Rogue"),
    -261: ("Class Quests", "Hunter"),
    -262: ("Class Quests", "Druid"),
    -263: ("Class Quests", "Priest"),
    
    # Profession quests (grouped under "Professions")
    -24: ("Professions", "Herbalism"),
    -101: ("Professions", "Fishing"),
    -121: ("Professions", "Blacksmithing"),
    -181: ("Professions", "Alchemy"),
    -182: ("Professions", "Leatherworking"),
    -201: ("Professions", "Engineering"),
    -264: ("Professions", "Tailoring"),
    -304: ("Professions", "Cooking"),
    -324: ("Professions", "First Aid"),
    -762: ("Professions", "Riding"),
    
    # Special quest types
    -1: ("Special", "Epic"),
    -21: ("Special", "Legendary"),
    -22: ("Special", "Seasonal"),
    -41: ("Special", "PvP"),
    -221: ("Special", "Treasure Map"),
    -241: ("Special", "Children's Week"),
    -284: ("Special", "Lunar Festival"),
    -364: ("Special", "Darkmoon Faire"),
    -365: ("Special", "Ahn'Qiraj War"),
    -366: ("Special", "Love is in the Air"),
    -367: ("Special", "Midsummer"),
    -368: ("Special", "Brewfest"),
    -369: ("Special", "Hallow's End"),
    -370: ("Special", "Pilgrim's Bounty"),
}

def main():
    print("Connecting to MySQL (tw_world)...")
    mysql_conn = mysql.connector.connect(**mysql_config)
    mysql_cursor = mysql_conn.cursor(dictionary=True)
    
    print("Connecting to MySQL (aowow) for zone names...")
    aowow_config = mysql_config.copy()
    aowow_config['database'] = 'aowow'
    aowow_conn = mysql.connector.connect(**aowow_config)
    aowow_cursor = aowow_conn.cursor(dictionary=True)
    
    print("Connecting to SQLite...")
    sqlite_conn = sqlite3.connect(SQLITE_DB)
    sqlite_cursor = sqlite_conn.cursor()
    
    # Get zones from aowow
    print("Fetching zones from aowow...")
    aowow_cursor.execute("""
        SELECT areatableID, name_loc0 
        FROM aowow_zones 
        WHERE name_loc0 IS NOT NULL AND name_loc0 != ''
    """)
    zones = {row['areatableID']: row['name_loc0'] for row in aowow_cursor.fetchall()}
    
    # Get quest counts per ZoneOrSort
    print("Fetching quest distribution...")
    mysql_cursor.execute("""
        SELECT ZoneOrSort, COUNT(*) as cnt 
        FROM quest_template 
        GROUP BY ZoneOrSort
    """)
    quest_counts = {row['ZoneOrSort']: row['cnt'] for row in mysql_cursor.fetchall()}
    
    # Create enhanced quest categories table
    print("Creating enhanced quest categories...")
    sqlite_cursor.executescript("""
        DROP TABLE IF EXISTS quest_category_groups;
        DROP TABLE IF EXISTS quest_categories_enhanced;
        
        CREATE TABLE quest_category_groups (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE
        );
        
        CREATE TABLE quest_categories_enhanced (
            id INTEGER PRIMARY KEY,
            group_id INTEGER NOT NULL,
            name TEXT NOT NULL,
            quest_count INTEGER DEFAULT 0,
            FOREIGN KEY (group_id) REFERENCES quest_category_groups(id)
        );
        
        CREATE INDEX idx_quest_cat_enhanced_group ON quest_categories_enhanced(group_id);
    """)
    
    # Insert groups
    groups = ["Zones", "Class Quests", "Professions", "Special", "Other"]
    group_ids = {}
    for group_name in groups:
        sqlite_cursor.execute(
            "INSERT INTO quest_category_groups (name) VALUES (?)",
            (group_name,)
        )
        group_ids[group_name] = sqlite_cursor.lastrowid
    
    # Process all ZoneOrSort values
    categories_added = 0
    for zone_sort, count in sorted(quest_counts.items(), key=lambda x: -x[1]):
        if zone_sort > 0:
            # Zone-based quest
            zone_name = zones.get(zone_sort, f"Zone {zone_sort}")
            group_id = group_ids["Zones"]
            sqlite_cursor.execute(
                "INSERT OR IGNORE INTO quest_categories_enhanced (id, group_id, name, quest_count) VALUES (?, ?, ?, ?)",
                (zone_sort, group_id, zone_name, count)
            )
        elif zone_sort < 0:
            # Special category
            if zone_sort in QUEST_SPECIAL_CATEGORIES:
                group_name, cat_name = QUEST_SPECIAL_CATEGORIES[zone_sort]
                group_id = group_ids[group_name]
            else:
                group_id = group_ids["Other"]
                cat_name = f"Category {abs(zone_sort)}"
            
            sqlite_cursor.execute(
                "INSERT OR IGNORE INTO quest_categories_enhanced (id, group_id, name, quest_count) VALUES (?, ?, ?, ?)",
                (zone_sort, group_id, cat_name, count)
            )
        categories_added += 1
    
    sqlite_conn.commit()
    
    # Verify
    sqlite_cursor.execute("SELECT COUNT(*) FROM quest_category_groups")
    group_count = sqlite_cursor.fetchone()[0]
    sqlite_cursor.execute("SELECT COUNT(*) FROM quest_categories_enhanced")
    cat_count = sqlite_cursor.fetchone()[0]
    
    print(f"\nDone! Created:")
    print(f"  - {group_count} category groups")
    print(f"  - {cat_count} quest categories")
    
    # Show summary
    print("\nCategory summary:")
    sqlite_cursor.execute("""
        SELECT g.name, COUNT(c.id) as cnt, SUM(c.quest_count) as quests
        FROM quest_category_groups g
        LEFT JOIN quest_categories_enhanced c ON g.id = c.group_id
        GROUP BY g.id
        ORDER BY quests DESC
    """)
    for row in sqlite_cursor.fetchall():
        print(f"  - {row[0]}: {row[1]} categories, {row[2] or 0} quests")
    
    mysql_cursor.close()
    mysql_conn.close()
    aowow_cursor.close()
    aowow_conn.close()
    sqlite_cursor.close()
    sqlite_conn.close()

if __name__ == "__main__":
    main()
