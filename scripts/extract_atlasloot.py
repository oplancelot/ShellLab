import os
import re
import json

# Base paths
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
ADDONS_DIR = os.path.join(BASE_DIR, '..', 'addons', 'AtlasLoot', 'Database')
DATA_DIR = os.path.join(BASE_DIR, '..', 'data')

CATEGORIES = [
    {"key": "AtlasLootInstances", "name": "Instances", "file": "Instances.lua", "sort": 1},
    {"key": "AtlasLootSets", "name": "Sets", "file": "Sets.lua", "sort": 2},
    {"key": "AtlasLootFactions", "name": "Factions", "file": "Factions.lua", "sort": 3},
    {"key": "AtlasLootPvP", "name": "PvP", "file": "PvP.lua", "sort": 4},
    {"key": "AtlasLootWorldBosses", "name": "World Bosses", "file": "WorldBosses.lua", "sort": 5},
    {"key": "AtlasLootWorldEvents", "name": "World Events", "file": "WorldEvents.lua", "sort": 6},
    {"key": "AtlasLootCrafting", "name": "Crafting", "file": "Crafting.lua", "sort": 7},
]

INSTANCE_PREFIXES = {
    "MC": "Molten Core",
    "Ony": "Onyxia's Lair",
    "BWL": "Blackwing Lair",
    "ZG": "Zul'Gurub",
    "AQ20": "Ruins of Ahn'Qiraj",
    "AQ40": "Temple of Ahn'Qiraj",
    "NAX": "Naxxramas",
    "BRD": "Blackrock Depths",
    "LBRS": "Lower Blackrock Spire",
    "UBRS": "Upper Blackrock Spire",
    "Strat": "Stratholme",
    "Scholo": "Scholomance",
    "DM": "Dire Maul",
    "ST": "Sunken Temple",
    "Mara": "Maraudon",
    "Uld": "Uldaman",
    "RFK": "Razorfen Kraul",
    "RFD": "Razorfen Downs",
    "SM": "Scarlet Monastery",
    "WC": "Wailing Caverns",
    "SFK": "Shadowfang Keep",
    "RFC": "Ragefire Chasm",
    "DM2": "Deadmines",
    "VC": "Deadmines",
}

def parse_table_register(path):
    display_names = {}
    if not os.path.exists(path):
        print(f"Warning: {path} not found")
        return display_names

    with open(path, 'r', encoding='utf-8', errors='ignore') as f:
        content = f.read()

    # Match entries like ["TableName"] = { ... }
    # This is simplified; the Go regex parsed line by line which is safer for this format
    # Let's use line processing
    
    key_regex = re.compile(r'\["(\w+)"\]\s*=')
    atlas_regex = re.compile(r'"AtlasLoot\w*Items"')
    al_regex = re.compile(r'AL\["([^"]+)"\]')
    quote_regex = re.compile(r'\{\s*"([^"]+)"')

    current_key = None
    buffer = []

    lines = content.splitlines()
    for line in lines:
        match = key_regex.search(line)
        if match:
            current_key = match.group(1)
            buffer = [line]
        elif current_key:
            buffer.append(line.strip())
            
            # Check for AtlasLoot marker
            combined = " ".join(buffer)
            if atlas_regex.search(combined):
                # Extract display name
                al_matches = al_regex.findall(combined)
                if al_matches:
                    parts = [p for p in al_matches if p not in ["Rare", "Summon", "Quest", "Enchants"]]
                    if parts:
                        display_names[current_key] = " - ".join(parts)
                else:
                    q_match = quote_regex.search(combined)
                    if q_match:
                         display_names[current_key] = q_match.group(1)
                
                current_key = None
                buffer = []
    
    return display_names

def clean_table_name(name):
    clean_name = name
    for prefix in ["MC", "BWL", "Ony", "ZG", "AQ20", "AQ40", "NAX", "BRD", "LBRS", "UBRS", "Strat", "Scholo", "DM", "ST", "Mara", "Uld"]:
        if clean_name.startswith(prefix):
            clean_name = clean_name[len(prefix):]
            break
            
    # Insert spaces (simplified)
    # Python regex for camelcase splitting is easier
    s1 = re.sub('(.)([A-Z][a-z]+)', r'\1 \2', clean_name)
    return re.sub('([a-z0-9])([A-Z])', r'\1 \2', s1).strip()

def get_display_name(table_name, display_names):
    if table_name in display_names:
        name = display_names[table_name]
    else:
        name = clean_table_name(table_name)
        
    # Append class names if missing (simplified logic from Go)
    classes = ["Druid", "Hunter", "Mage", "Paladin", "Priest", "Rogue", "Shaman", "Warlock", "Warrior"]
    if "Class Set" in name or "Tier" in name:
        for cls in classes:
            if cls in table_name and cls not in name:
                name = f"{cls} {name}"
                break
    else:
         for cls in classes:
            if cls in table_name and cls not in name:
                name = f"{cls} {name}"
                break

    if name.endswith(" C"):
        name = name[:-2] + " (Compact)"
    
    if table_name.startswith("DME") and "East" not in name: name += " (East)"
    elif table_name.startswith("DMN") and "North" not in name: name += " (North)"
    elif table_name.startswith("DMW") and "West" not in name: name += " (West)"
    
    return name

def parse_lua_file(path):
    tables = {}
    if not os.path.exists(path):
        return tables

    with open(path, 'r', encoding='utf-8', errors='ignore') as f:
        content = f.read()

    # Regex patterns
    table_start_regex = re.compile(r'^\s*(\w+)\s*=\s*\{')
    item_regex = re.compile(r'\{\s*(\d+)\s*,.*?"([^"]*%)')
    simple_item_regex = re.compile(r'\{\s*(\d+)\s*,')

    current_table = None
    current_items = []

    for line in content.splitlines():
        # Table start
        match = table_start_regex.search(line)
        if match:
            if current_table and current_items:
                tables[current_table] = {"name": current_table, "items": current_items}
            current_table = match.group(1)
            current_items = []
            continue

        # Item entry
        if "{" in line and current_table:
            match = item_regex.search(line)
            if match:
                item_id = int(match.group(1))
                drop_rate = match.group(2)
                current_items.append({"id": item_id, "drop_rate": drop_rate})
            else:
                match = simple_item_regex.search(line)
                if match:
                    item_id = int(match.group(1))
                    current_items.append({"id": item_id, "drop_rate": ""})

        # Table end
        if "};" in line and current_table:
            if current_items:
                tables[current_table] = {"name": current_table, "items": current_items}
            current_table = None
            current_items = []

    # Filter Compact tables
    filtered = {}
    for name, table in tables.items():
        if name.endswith("C"):
            base = name[:-1]
            if base in tables:
                continue
        filtered[name] = table
    
    return filtered

def group_tables(tables, category_key):
    if category_key != "AtlasLootInstances":
        # Group all in one module named after category
        # But maybe we want Modules? 
        # For Sets/Factions, usually they are just lists. 
        # The Go code created a single module with the same name as category.
        return {"Default": list(tables.keys())}

    groups = {"Other": []}
    for name in tables:
        matched = False
        for prefix, instance in INSTANCE_PREFIXES.items():
            if name.startswith(prefix):
                if instance not in groups: groups[instance] = []
                groups[instance].append(name)
                matched = True
                break
        if not matched:
            groups["Other"].append(name)
    
    # Remove Other if empty
    if not groups["Other"]:
        del groups["Other"]
    return groups

def main():
    print("Extracting AtlasLoot data...")
    os.makedirs(DATA_DIR, exist_ok=True)

    # 1. Parse Register
    register_path = os.path.join(ADDONS_DIR, "TableRegister.lua")
    display_names = parse_table_register(register_path)
    print(f"Loaded {len(display_names)} display name mappings.")

    # 2. Process Categories
    result_data = []

    for cat in CATEGORIES:
        print(f"Processing {cat['name']}...")
        lua_path = os.path.join(ADDONS_DIR, cat['file'])
        tables = parse_lua_file(lua_path)
        
        module_groups = group_tables(tables, cat['key'])
        
        modules_list = []
        for mod_name, table_names in module_groups.items():
            mod_data = {
                "key": mod_name if mod_name != "Default" else cat['name'],
                "name": mod_name if mod_name != "Default" else cat['name'],
                "tables": []
            }
            
            for t_name in table_names:
                table_def = tables[t_name]
                display_name = get_display_name(t_name, display_names)
                
                mod_data["tables"].append({
                    "key": t_name,
                    "name": display_name,
                    "items": table_def["items"]
                })
            
            # Sort tables by name? Or keep order?
            # Go code didn't sort explicitly, just used iteration order (random in Go map!)
            # Python dicts preserve insertion order in new versions, but `group_tables` iterated keys.
            modules_list.append(mod_data)

        result_data.append({
            "key": cat['key'],
            "name": cat['name'],
            "sort": cat['sort'],
            "modules": modules_list
        })

    # 3. Save to JSON
    output_path = os.path.join(DATA_DIR, "atlasloot.json")
    with open(output_path, 'w', encoding='utf-8') as f:
        json.dump(result_data, f, indent=2, ensure_ascii=False)
    
    print(f"Successfully saved to {output_path}")

if __name__ == "__main__":
    main()
