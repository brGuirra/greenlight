INSERT INTO movies (title, year, runtime, genres)
VALUES
('Inception', 2010, 148, ARRAY['Sci-Fi', 'Action', 'Thriller']),
('The Shawshank Redemption', 1994, 142, ARRAY['Drama']),
('The Godfather', 1972, 175, ARRAY['Crime', 'Drama']),
('The Dark Knight', 2008, 152, ARRAY['Action', 'Crime', 'Drama']),
('Pulp Fiction', 1994, 154, ARRAY['Crime', 'Drama']),
('Forrest Gump', 1994, 142, ARRAY['Drama', 'Romance']),
('The Matrix', 1999, 136, ARRAY['Action', 'Sci-Fi']),
('Schindler''s List', 1993, 195, ARRAY['Biography', 'Drama', 'History']),
(
    'The Lord of the Rings: The Fellowship of the Ring',
    2001,
    178,
    ARRAY['Action', 'Adventure', 'Drama']
),
('Avatar', 2009, 162, ARRAY['Action', 'Adventure', 'Fantasy']);
