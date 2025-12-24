import db_config
import mysql.connector
import sqlite3
import os

mysql_config = db_config.get_mysql_config()
tw_config = {**db_config.get_mysql_config(db_config.get_tw_world_db()), 'charset': 'utf8mb4'}
SQLITE_DB = os.path.join(os.path.dirname(__file__), '..', 'data', 'shelllab.db')

# Standard Vanilla Quest Sorts
QUEST_SORTS = {
    -1: "Epic",
    -22: "Seasonal",
    -24: "Herbalism",
    -25: "Battlegrounds",
    -44: "Other",
    -61: "Warlock",
    -81: "Warrior",
    -82: "Shaman",
    -101: "Fishing",
    -121: "Blacksmithing",
    -141: "Paladin",
    -161: "Mage",
    -162: "Rogue",
    -181: "Alchemy",
    -182: "Leatherworking",
    -201: "Engineering",
    -221: "Treasure Map",
    -241: "Tailoring",
    -261: "Hunter",
    -262: "Priest",
    -263: "Druid",
    -264: "First Aid",
    -284: "Special",
    -304: "Cooking",
    -324: "Other",
    -344: "Legendaries",
    -364: "Darkmoon Faire",
    -365: "Ahn'Qiraj War Effort",
    -366: "Lunar Festival",
    -367: "Reputation",
    -368: "Invasion",
    -369: "Midsummer",
    -370: "Brewfest",
    # Turtle WoW Custom (Guessed or to be verified)
    -1001: "Survival",
    -1002: "Gardening"
}

def export_categories():
    print("Connecting to MySQL...")
    try:
        # 1. Fetch Zones from aowow
        mysql_conn = mysql.connector.connect(**mysql_config)
        mysql_cursor = mysql_conn.cursor(dictionary=True)
        
        print("Fetching zones from aowow_zones...")
        mysql_cursor.execute("SELECT areatableID, name_loc0 FROM aowow_zones")
        zones = mysql_cursor.fetchall()
        print(f"Found {len(zones)} zones")
        
        mysql_cursor.close()
        mysql_conn.close()

        # 2. Fetch distinct negative ZoneOrSort from quest_template to find any missing sorts
        print("Checking for custom sorts in quest_template...")
        tw_conn = mysql.connector.connect(**tw_config)
        tw_cursor = tw_conn.cursor()
        tw_cursor.execute("SELECT DISTINCT ZoneOrSort FROM quest_template WHERE ZoneOrSort < 0")
        used_sorts = [row[0] for row in tw_cursor.fetchall()]
        tw_cursor.close()
        tw_conn.close()

        print(f"Found {len(used_sorts)} used sort IDs.")
        for sort_id in used_sorts:
            if sort_id not in QUEST_SORTS:
                print(f"Warning: Unknown Sort ID {sort_id}, adding as 'Unknown Sort {sort_id}'")
                QUEST_SORTS[sort_id] = f"Unknown Sort {sort_id}"

        # 3. Write to SQLite
        print(f"Connecting to SQLite: {SQLITE_DB}")
        sqlite_conn = sqlite3.connect(SQLITE_DB)
        c = sqlite_conn.cursor()
        
        c.execute("DROP TABLE IF EXISTS quest_categories")
        c.execute("CREATE TABLE quest_categories (id INTEGER PRIMARY KEY, name TEXT)")
        
        print("Inserting zones (positive IDs)...")
        for z in zones:
            c.execute("INSERT OR REPLACE INTO quest_categories VALUES (?,?)", (z['areatableID'], z['name_loc0']))
            
        print("Inserting sorts (negative IDs)...")
        for sort_id, name in QUEST_SORTS.items():
            c.execute("INSERT OR REPLACE INTO quest_categories VALUES (?,?)", (sort_id, name))
        
        sqlite_conn.commit()
        sqlite_conn.close()
        print("Done!")
    except Exception as e:
        print(f"Error: {e}")

if __name__ == '__main__':
    export_categories()
