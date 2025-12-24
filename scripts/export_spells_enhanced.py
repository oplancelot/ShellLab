import os
import json
import mysql.connector
import db_config
from mysql.connector import Error

DATA_DIR = os.path.join(os.path.dirname(os.path.dirname(__file__)), 'data')

def export_spells():
    conn = None
    try:
        print("Connecting to MySQL...")
        conn = mysql.connector.connect(**db_config.get_mysql_config())
        cursor = conn.cursor(dictionary=True)

        # 1. Find AoWoW DB
        aowow_db = None
        cursor.execute("SHOW DATABASES LIKE '%aowow%'")
        dbs = cursor.fetchall()
        if dbs:
            aowow_db = list(dbs[0].values())[0]
            print(f"Found AoWoW DB: {aowow_db}")

        # 2. Fetch Duration Data (from AoWoW if avail)
        if aowow_db:
            try:
                print("Fetching Spell Durations from AoWoW...")
                cursor.execute(f"SHOW COLUMNS FROM {aowow_db}.aowow_spellduration")
                cols = [c['Field'] for c in cursor.fetchall()]
                
                col_id = 'durationID' if 'durationID' in cols else 'id'
                col_base = 'durationBase'
                col_per = 'durationPerLevel' if 'durationPerLevel' in cols else '0'
                col_max = 'maxDuration' if 'maxDuration' in cols else '0'
                
                query = f"SELECT {col_id} as id, {col_base} as durationBase, {col_per} as durationPerLevel, {col_max} as maxDuration FROM {aowow_db}.aowow_spellduration"
                cursor.execute(query)
                durations = cursor.fetchall()
                
                with open(os.path.join(DATA_DIR, 'spell_durations.json'), 'w', encoding='utf-8') as f:
                    json.dump(durations, f, ensure_ascii=False)
                print(f"Saved {len(durations)} durations.")
            except Exception as e:
                print(f"Error fetching durations: {e}")

        # 3. Fetch Spell Icons (from AoWoW if avail)
        icon_map = {}
        if aowow_db:
            try:
                print("Fetching Spell Icons from AoWoW...")
                cursor.execute(f"SELECT id, iconname FROM {aowow_db}.aowow_spellicons")
                for r in cursor.fetchall():
                    icon_map[r['id']] = r['iconname']
                print(f"Loaded {len(icon_map)} icons.")
            except Exception as e:
                print(f"Error fetching icons: {e}")

        # 4. Fetch Spells (from TW_WORLD)
        print("Fetching Spells from World DB...")
        table_name = "spell_template"
        
        cursor.execute(f"SHOW COLUMNS FROM {table_name}")
        cols = {c['Field'] for c in cursor.fetchall()}
        
        select_parts = []
        select_parts.append("entry")
        select_parts.append("name")
        select_parts.append("description")
        
        # Effects
        for i in range(1, 4):
            base = f"effectBasePoints{i}"
            if f"effect_base_points{i}" in cols:
                select_parts.append(f"effect_base_points{i} as effectBasePoints{i}")
            elif base in cols:
                select_parts.append(base)
            else:
                select_parts.append(f"0 as effectBasePoints{i}")

            sides = f"effectDieSides{i}"
            if f"effect_die_sides{i}" in cols:
                select_parts.append(f"effect_die_sides{i} as effectDieSides{i}")
            elif sides in cols:
                select_parts.append(sides)
            else:
                 select_parts.append(f"0 as effectDieSides{i}")

        # Duration
        if "durationIndex" in cols:
            select_parts.append("durationIndex")
        elif "duration_index" in cols:
             select_parts.append("duration_index as durationIndex")
        else:
            select_parts.append("0 as durationIndex")

        # Icon
        if "spellIconId" in cols:
             select_parts.append("spellIconId")
        elif "activeIconId" in cols:
             select_parts.append("activeIconId as spellIconId")
        else:
             select_parts.append("0 as spellIconId")
             
        query = f"SELECT {', '.join(select_parts)} FROM {table_name}"
        cursor.execute(query)
        
        spells = []
        for row in cursor.fetchall():
            # Map Icon
            icon_id = row.get('spellIconId', 0)
            row['iconName'] = icon_map.get(icon_id, '')
            if row['iconName']:
                 row['iconName'] = row['iconName'].lower()

            # Clean
            if row['description'] is None: row['description'] = ""
            if row['name'] is None: row['name'] = ""
            
            spells.append(row)
            
        print(f"Fetched {len(spells)} spells.")
        with open(os.path.join(DATA_DIR, 'spells_enhanced.json'), 'w', encoding='utf-8') as f:
            json.dump(spells, f, ensure_ascii=False)

    except Error as e:
        print(f"Error: {e}")
    finally:
        if conn and conn.is_connected():
            cursor.close()
            conn.close()

if __name__ == "__main__":
    export_spells()
