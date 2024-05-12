-- +goose Up
-- +goose StatementBegin
INSERT INTO "spender"(name, email)
VALUES
('John Doe', 'john_d@test.com'),
('Sara Doe', 'sara_d@test.com');

INSERT INTO "transaction"(date, amount, category, transaction_type, note, image_url, spender_id)
VALUES
('2024-05-12', 200, 'Salary', 'income', 'Monthly salary', 'url_to_image1', 1),
('2024-05-12', 150, 'Utilities', 'expense', 'Electricity bill payment', 'url_to_image2', 1),
('2024-05-12', 500, 'Freelancing', 'income', 'Freelance web design payment', 'url_to_image3', 1),
('2024-05-12', 120, 'Entertainment', 'expense', 'Movie tickets for family', 'url_to_image4', 1),
('2024-05-12', 300, 'Bonus', 'income', 'Annual bonus', 'url_to_image5', 1),
('2024-05-12', 250, 'Travel', 'expense', 'Train tickets booking', 'url_to_image6', 1),
('2024-05-12', 400, 'Dividends', 'income', 'Stock dividends', 'url_to_image7', 1),
('2024-05-12', 45, 'Stationery', 'expense', 'Office supplies', 'url_to_image8', 1),
('2024-05-12', 150, 'Rental income', 'income', 'Monthly rent from property', 'url_to_image9', 1),
('2024-05-12', 35, 'Books', 'expense', 'Latest novel purchase', 'url_to_image10', 1),
('2024-05-12', 250, 'Consulting', 'income', 'Consulting fees received', 'url_to_image11', 1),
('2024-05-12', 60, 'Dining', 'expense', 'Lunch at a restaurant', 'url_to_image12', 1),
('2024-05-12', 220, 'Investment Return', 'income', 'Return on investments', 'url_to_image13', 1),
('2024-05-12', 190, 'Furniture', 'expense', 'New office chair', 'url_to_image14', 1),
('2024-05-12', 100, 'Interest', 'income', 'Bank account interest', 'url_to_image15', 1),
('2024-05-12', 110, 'Gardening', 'expense', 'Garden tools and plants', 'url_to_image16', 1),
('2024-05-12', 180, 'Sales', 'income', 'Sales from online store', 'url_to_image17', 1),
('2024-05-12', 250, 'Gifts', 'expense', 'Birthday gifts for family', 'url_to_image18', 1),
('2024-05-12', 130, 'Commissions', 'income', 'Sales commission received', 'url_to_image19', 1),
('2024-05-12', 200, 'Insurance', 'expense', 'Car insurance renewal', 'url_to_image20', 1),
('2024-05-12', 230, 'Royalties', 'income', 'Book royalties', 'url_to_image21', 1),
('2024-05-12', 180, 'Technology', 'expense', 'New laptop', 'url_to_image22', 1),
('2024-05-12', 90, 'Transportation', 'expense', 'Taxi fares', 'url_to_image23', 1),
('2024-05-12', 85, 'Health', 'expense', 'Dental checkup', 'url_to_image24', 1),
('2024-05-12', 95, 'Fitness', 'income', 'Personal training sessions', 'url_to_image25', 1),
('2024-05-12', 105, 'Pets', 'expense', 'Veterinary visit', 'url_to_image26', 1),
('2024-05-12', 140, 'Stationery', 'income', 'Sold art supplies online', 'url_to_image27', 2),
('2024-05-12', 260, 'Electronics', 'expense', 'New camera', 'url_to_image28', 2),
('2024-05-12', 125, 'Books', 'income', 'Used book sales', 'url_to_image29', 2),
('2024-05-12', 180, 'Clothing', 'expense', 'Winter jacket and boots', 'url_to_image30', 2),
('2024-05-12', 70, 'Dining', 'income', 'Catering service payment', 'url_to_image31', 2),
('2024-05-12', 45, 'Sports', 'expense', 'Golf club membership', 'url_to_image32', 2),
('2024-05-12', 210, 'Furniture', 'income', 'Old furniture sold', 'url_to_image33', 2),
('2024-05-12', 330, 'Renovation', 'expense', 'Kitchen remodeling', 'url_to_image34', 2),
('2024-05-12', 120, 'Gardening', 'income', 'Landscaping services provided', 'url_to_image35', 2),
('2024-05-12', 65, 'Beauty', 'expense', 'Makeup and skincare products', 'url_to_image36', 2),
('2024-05-12', 270, 'Gifts', 'income', 'Handmade gifts sold', 'url_to_image37', 2),
('2024-05-12', 150, 'Hobbies', 'expense', 'Model building kits', 'url_to_image38', 2),
('2024-05-12', 230, 'Insurance', 'income', 'Life insurance policy matured', 'url_to_image39', 2),
('2024-05-12', 190, 'Technology', 'expense', 'Smart home devices', 'url_to_image40', 2),
('2024-05-12', 240, 'Education', 'income', 'Tutoring services', 'url_to_image41', 2),
('2024-05-12', 100, 'Transportation', 'expense', 'Bike repair', 'url_to_image42', 2),
('2024-05-12', 95, 'Health', 'income', 'Health workshop conducted', 'url_to_image43', 2),
('2024-05-12', 85, 'Fitness', 'expense', 'New gym equipment', 'url_to_image44', 2),
('2024-05-12', 115, 'Pets', 'income', 'Pet grooming services', 'url_to_image45', 2),
('2024-05-12', 160, 'Stationery', 'expense', 'High-end art materials', 'url_to_image46', 2),
('2024-05-12', 290, 'Electronics', 'income', 'Old electronics sold', 'url_to_image47', 2),
('2024-05-12', 135, 'Books', 'expense', 'Subscription to literary journals', 'url_to_image48', 2),
('2024-05-12', 200, 'Clothing', 'income', 'Vintage clothes sold', 'url_to_image49', 2),
('2024-05-12', 90, 'Dining', 'expense', 'Organic food supplies', 'url_to_image50', 2);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "spender"
WHERE email in ('john_d@test.com', 'sara_d@test.com');

DELETE FROM "transaction"
WHERE
    spender_id in (
        1,2
    );
-- +goose StatementEnd
