
import os
import requests
from bs4 import BeautifulSoup
import uuid
import re

# Categories to scrape
CATEGORIES = [
    {"name": "Бургери", "slug": "burgers"},
    {"name": "Сети", "slug": "sets"},
    {"name": "Перші страви", "slug": "persi stravi"},
    {"name": "Соуси", "slug": "sauces"},
    {"name": "Напої", "slug": "drinks"},
    {"name": "Фрі та сир", "slug": "fries-and-cheese"},
    {"name": "Салати", "slug": "salad"},
    {"name": "Гарячі страви", "slug": "hot"},
    {"name": "Гарячі закуски", "slug": "snacks"},
    {"name": "Холодні закуски", "slug": "holodni"},
]

BASE_URL = "https://brobar.delivery"
UPLOADS_DIR = "uploads"

# Ensure uploads directory exists
if not os.path.exists(UPLOADS_DIR):
    os.makedirs(UPLOADS_DIR)

def download_image(url, filename):
    try:
        response = requests.get(url, stream=True)
        if response.status_code == 200:
            filepath = os.path.join(UPLOADS_DIR, filename)
            with open(filepath, 'wb') as f:
                for chunk in response.iter_content(1024):
                    f.write(chunk)
            return filepath
    except Exception as e:
        print(f"Failed to download {url}: {e}")
    return None

def main():
    # Placeholder for logic
    print("Ready to scrape")

if __name__ == "__main__":
    main()
