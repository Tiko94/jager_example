INSERT INTO authors (id, first_name, last_name, username) VALUES
  (1, "John", "Cleese", "jcleese"),
  (2, "Terry", "Jones", "tjones"),
  (3, "Eric", "Idle", "eidle"),
  (4, "Terry", "Gilliam", "tgilliam"),
  (5, "Graham", "Chapman", "gchapman"),
  (6, "Michael", "Palin", "mpalin");

INSERT INTO posts (subject, body, author_id) VALUES
  ("From Killer Rabbits to African Swallows - How to Find an Unusual Pet", "", 1),
  ("Ducks and Witches - Why Physics Matters", "", 2),
  ("Why I Think Sir Robin is the Bravest Man I Know", "", 3),
  ("Why I Took the Limb Dismemberment Class", "", 5),
  ("Swamp Castle - A Study in Modern Architecure Gone Wrong", "", 6);
