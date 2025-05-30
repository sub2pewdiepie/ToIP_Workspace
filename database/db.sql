-- Таблица учебных групп
CREATE TABLE Academic_Groups (
    academic_group_id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица групп
CREATE TABLE Groups (
    group_id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    academic_group_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (academic_group_id) REFERENCES Academic_Groups(academic_group_id)
);

-- Таблица пользователей
CREATE TABLE Users (
    user_id INTEGER PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Связующая таблица для групп и пользователей
CREATE TABLE Group_Users (
    group_id INTEGER,
    user_id INTEGER,
    role VARCHAR(50) DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES Groups(group_id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);

-- Таблица предметов
CREATE TABLE Subjects (
    subject_id INTEGER PRIMARY KEY,
    group_id INTEGER,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES Groups(group_id)
);

-- Таблица заданий


-- Связующая таблица для заданий и пользователей
CREATE TABLE Task_Users (
    task_id INTEGER,
    user_id INTEGER,
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP,
    PRIMARY KEY (task_id, user_id),
    FOREIGN KEY (task_id) REFERENCES Tasks(task_id),
    FOREIGN KEY (user_id) REFERENCES Users(user_id)
);

-- Таблица материалов
CREATE TABLE Materials (
    material_id INTEGER PRIMARY KEY,
    subject_id INTEGER,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    created_by INTEGER,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (subject_id) REFERENCES Subjects(subject_id),
    FOREIGN KEY (created_by) REFERENCES Users(user_id)
);

-- Таблица временных слотов (пары)
CREATE TABLE Time_Slots (
    slot_id INTEGER PRIMARY KEY,
    slot_number INTEGER NOT NULL CHECK (slot_number BETWEEN 1 AND 9),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    UNIQUE (slot_number)
);

-- Таблица расписания
CREATE TABLE Schedules (
    schedule_id INTEGER PRIMARY KEY,
    group_id INTEGER,
    subject_id INTEGER,
    teacher_initials VARCHAR(50),
    classroom VARCHAR(50),
    slot_id INTEGER,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES Groups(group_id),
    FOREIGN KEY (subject_id) REFERENCES Subjects(subject_id),
    FOREIGN KEY (slot_id) REFERENCES Time_Slots(slot_id)
);