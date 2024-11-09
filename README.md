# Recipe Finder App

## Description

Recipe Finder is an application that helps find recipes based on provided ingredients. Users can search for recipes by ingredients, add their own recipes, save favorites, and also edit or delete recipes they have created.

## Key Features

- **Recipe Search by Ingredients**: Find recipes that can be made with available ingredients.
- **Recipe Creation**: Add your own recipes with ingredients and instructions.
- **Save Recipes**: Add recipes to favorites for quick access.
- **Edit and Delete**: Manage your recipes—update them or delete if they are no longer needed.

## Tech Stack

- **Backend**: Go (Gin)
- **Database**: PostgreSQL
- **Cache**: Redis

## Launch
# 1
git clone https://github.com/Perseverance7/recipes.git

# 2
docker-compose up --build

# 3
if app run{
    Swagger Api: http://localhost:8081/swagger/index.html
} 
else {
    docs/swagger.yaml
}

