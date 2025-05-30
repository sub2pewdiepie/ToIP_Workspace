-- Выполнится только при первом запуске контейнера
DO
$$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_database WHERE datname = 'mydb'
   ) THEN
      CREATE DATABASE mydb;
   END IF;
END
$$;
