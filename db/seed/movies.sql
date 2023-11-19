INSERT INTO movies (title, year, runtime, genres)
VALUES
('Whiplash', 1921, 122, '{"Sci-Fi"}'),
('The Dark Knight', 1955, 145, '{"Romance"}'),
('The Green Mile', 1906, 102, '{"Fantasy"}'),
('Gladiator', 1952, 174, '{"Family"}'),
('Forrest Gump', 2016, 108, '{"History"}'),
('Harry Potter and the Deathly Hallows - Part 2', 1984, 139, '{"Crime"}'),
('Inception', 1946, 118, '{"Family"}'),
('Rocky', 1945, 101, '{"Musical"}'),
('Interstellar', 1965, 146, '{"Mystery"}'),
('Green Book', 1937, 147, '{"Musical"}'),
('The Usual Suspects', 1932, 176, '{"Western"}'),
('3 Idiots', 1953, 158, '{"Music"}'),
('Prisoners', 1931, 130, '{"Crime"}'),
('The Terminator', 1994, 136, '{"Animation"}'),
('WALL·E', 1936, 156, '{"Fantasy"}'),
('Casablanca', 1921, 122, '{"Biography"}'),
('Léon', 1912, 113, '{"Action"}'),
('V for Vendetta', 1913, 129, '{"Action"}'),
('Apocalypse Now', 1953, 123, '{"Western"}'),
('Joker', 1992, 113, '{"Family"}'),
('Saving Private Ryan', 1942, 105, '{"War"}'),
('Monsters, Inc.', 1952, 128, '{"Biography"}'),
('La vita è bella', 1924, 102, '{"Musical"}'),
('Memento', 1964, 108, '{"War"}'),
('The Lion King', 1989, 134, '{"Fantasy"}'),
('Monty Python and the Holy Grail', 1928, 171, '{"Romance"}'),
('Top Gun: Maverick', 2016, 142, '{"Action"}'),
('Terminator 2: Judgment Day', 2000, 154, '{"Biography"}'),
('The Shawshank Redemption', 1922, 99, '{"Biography"}'),
('The Pianist', 1963, 95, '{"Sci-Fi"}'),
('One Flew Over the Cuckoo''s Nest', 1992, 104, '{"Biography"}'),
('Die Hard', 2013, 172, '{"Biography"}'),
('Snatch', 1954, 108, '{"Sport"}'),
('Gandhi', 2007, 109, '{"Biography"}'),
('Judgment at Nuremberg', 2000, 121, '{"Romance"}'),
('The Matrix', 1953, 172, '{"Sci-Fi", "Adventure"}'),
('Indiana Jones and the Last Crusade', 2000, 160, '{"Film-Noir"}'),
('Once Upon a Time in America', 1921, 142, '{"Film-Noir"}'),
('Gone with the Wind', 1993, 133, '{"History"}'),
('Grease', 1978, 110, '{"Comedy"}'),
('Dirty Dancing', 1987, 100, '{"Romance"}'),
('Home Alone', 1990, 103, '{"Comedy", "Adventure"}'),
('Love Actually', 2003, 135, '{"Romande"}'),
('Titanic', 1997, 194, '{"Romance"}'),
('Jurassic Park', 1993, 127, '{"Action"}'),
('Mary Poppins', 1964, 139, '{"Comedy"}'),
('Toy Story', 1995, 81, '{"Comedy", "Adventure", "Animation"}'),
('Jaws', 1975, 124, '{"Thiller"}'),
('Top Gun', 1986, 109, '{"Action"}'),
('The Silence of the Lambs', 1991, 118, '{"Thiller"}')
RETURNING id, created_at, title, year, runtime, genres, version;
