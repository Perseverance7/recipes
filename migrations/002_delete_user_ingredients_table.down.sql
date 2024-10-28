-- Создание таблицы user_ingredients
CREATE TABLE user_ingredients (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    ingredient_id INT REFERENCES ingredients(id) ON DELETE CASCADE,
    quantity DECIMAL(10, 2) NOT NULL,
    PRIMARY KEY (user_id, ingredient_id)
);
