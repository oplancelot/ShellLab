#!/usr/bin/env python3
"""
Fix missing item icons by fetching from Turtle WoW database website
and downloading the icon images.
"""

import sqlite3
import requests
import re
import os
import time
from pathlib import Path
from bs4 import BeautifulSoup

# Configuration
DB_PATH = "../data/shelllab.db"
ICON_DIR = "../frontend/public/items/icons"
BASE_URL = "https://database.turtlecraft.gg/?item="
DELAY_SECONDS = 0.5  # Be nice to the server

def get_missing_icons():
    """Get list of items with missing icon_path"""
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()
    
    cursor.execute("""
        SELECT entry, name 
        FROM item_template 
        WHERE icon_path IS NULL OR icon_path = ''
        ORDER BY entry
    """)
    
    missing = cursor.fetchall()
    conn.close()
    
    print(f"Found {len(missing)} items with missing icons")
    return missing

def fetch_icon_from_website(entry):
    """Fetch icon name from Turtle WoW database website"""
    url = f"{BASE_URL}{entry}"
    
    try:
        response = requests.get(url, timeout=10)
        response.raise_for_status()
        
        # Parse HTML to find icon
        soup = BeautifulSoup(response.text, 'html.parser')
        
        # Look for: <div class="iconlarge" style="background-image: url(&quot;images/icons/large/inv_pants_cloth_03.png&quot;);">
        icon_div = soup.find('div', class_='iconlarge')
        
        if icon_div and icon_div.get('style'):
            style = icon_div['style']
            # Extract icon filename from style
            match = re.search(r'images/icons/large/([^.]+)\.png', style)
            if match:
                icon_name = match.group(1)
                print(f"  Found icon: {icon_name}")
                return icon_name
        
        print(f"  No icon found in HTML")
        return None
        
    except Exception as e:
        print(f"  Error fetching: {e}")
        return None

def update_icon_path(entry, icon_name):
    """Update icon_path in database"""
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()
    
    cursor.execute("""
        UPDATE item_template 
        SET icon_path = ? 
        WHERE entry = ?
    """, (icon_name, entry))
    
    conn.commit()
    conn.close()

def download_icon(icon_name):
    """Download icon image if not already present"""
    icon_path = Path(ICON_DIR) / f"{icon_name.lower()}.jpg"
    
    if icon_path.exists():
        return True
    
    # Try local cache first from Turtle WoW
    url = f"https://database.turtlecraft.gg/images/icons/large/{icon_name}.png"
    
    try:
        response = requests.get(url, timeout=10)
        if response.status_code == 200:
            # Convert PNG to JPG and save
            icon_path.parent.mkdir(parents=True, exist_ok=True)
            with open(icon_path, 'wb') as f:
                f.write(response.content)
            print(f"    Downloaded: {icon_name}.jpg")
            return True
    except:
        pass
    
    # Fallback to Wowhead CDN
    url = f"https://wow.zamimg.com/images/wow/icons/medium/{icon_name.lower()}.jpg"
    try:
        response = requests.get(url, timeout=10)
        if response.status_code == 200:
            icon_path.parent.mkdir(parents=True, exist_ok=True)
            with open(icon_path, 'wb') as f:
                f.write(response.content)
            print(f"    Downloaded from CDN: {icon_name}.jpg")
            return True
    except:
        pass
    
    print(f"    Failed to download: {icon_name}")
    return False

def main():
    print("=== Fix Missing Item Icons ===\n")
    
    # Get missing icons
    missing_items = get_missing_icons()
    
    if not missing_items:
        print("No missing icons!")
        return
    
    # Ask for confirmation
    print(f"\nThis will:")
    print(f"  1. Fetch icon data for {len(missing_items)} items from {BASE_URL}")
    print(f"  2. Update database with icon names")
    print(f"  3. Download missing icon images")
    print(f"\nDelay between requests: {DELAY_SECONDS}s")
    
    response = input("\nProceed? (yes/no): ")
    if response.lower() not in ['yes', 'y']:
        print("Cancelled.")
        return
    
    # Process each item
    success_count = 0
    download_count = 0
    
    for i, (entry, name) in enumerate(missing_items, 1):
        print(f"\n[{i}/{len(missing_items)}] Item {entry}: {name}")
        
        # Fetch icon name
        icon_name = fetch_icon_from_website(entry)
        
        if icon_name:
            # Update database
            update_icon_path(entry, icon_name)
            success_count += 1
            
            # Download icon
            if download_icon(icon_name):
                download_count += 1
        
        # Be nice to the server
        if i < len(missing_items):
            time.sleep(DELAY_SECONDS)
    
    print(f"\n=== Summary ===")
    print(f"Icons fetched: {success_count}/{len(missing_items)}")
    print(f"Icons downloaded: {download_count}")
    print(f"\nDatabase updated: {DB_PATH}")
    print(f"Icons saved to: {ICON_DIR}")

if __name__ == "__main__":
    main()
