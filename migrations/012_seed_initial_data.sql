-- Seed data for ACTIVITY_CATEGORIES
INSERT INTO activity_categories (id, name, description, icon, color) VALUES
    (uuid_generate_v4(), 'Olahraga', 'Kegiatan olahraga dan fitness', '🏃', '#4CAF50'),
    (uuid_generate_v4(), 'Makan', 'Waktu makan dan nutrisi', '🍽️', '#FF9800'),
    (uuid_generate_v4(), 'Kerja', 'Kegiatan kerja dan meeting', '💼', '#2196F3'),
    (uuid_generate_v4(), 'Istirahat', 'Waktu istirahat dan relaksasi', '😴', '#9C27B0'),
    (uuid_generate_v4(), 'Belajar', 'Kegiatan belajar dan pengembangan diri', '📚', '#F44336'),
    (uuid_generate_v4(), 'Hobi', 'Kegiatan hobi dan rekreasi', '🎨', '#E91E63'),
    (uuid_generate_v4(), 'Kesehatan', 'Kegiatan kesehatan dan medis', '🏥', '#00BCD4'),
    (uuid_generate_v4(), 'Sosial', 'Kegiatan sosial dan pertemuan', '👥', '#FFC107'),
    (uuid_generate_v4(), 'Rumah Tangga', 'Kegiatan rumah tangga dan perawatan', '🏠', '#795548'),
    (uuid_generate_v4(), 'Lainnya', 'Kegiatan lainnya', '📝', '#607D8B')
ON CONFLICT (name) DO NOTHING;

-- Seed data for RECOMMENDATION_TYPES
INSERT INTO recommendation_types (id, name, description, trigger_condition) VALUES
    (uuid_generate_v4(), 'exercise', 'Rekomendasi olahraga dan aktivitas fisik', 'Kurang aktivitas fisik atau pola olahraga tidak teratur'),
    (uuid_generate_v4(), 'nutrition', 'Rekomendasi nutrisi dan pola makan', 'Pola makan tidak teratur atau kurang nutrisi'),
    (uuid_generate_v4(), 'sleep', 'Rekomendasi tidur dan istirahat', 'Kurang tidur atau pola tidur tidak teratur'),
    (uuid_generate_v4(), 'hydration', 'Rekomendasi hidrasi dan minum air', 'Kurang konsumsi air atau dehidrasi'),
    (uuid_generate_v4(), 'stress', 'Rekomendasi manajemen stress', 'Tingkat stress tinggi atau terlalu banyak kegiatan'),
    (uuid_generate_v4(), 'work_life_balance', 'Rekomendasi keseimbangan kerja dan hidup', 'Terlalu banyak kegiatan kerja atau kurang waktu istirahat'),
    (uuid_generate_v4(), 'activity_variety', 'Rekomendasi variasi kegiatan', 'Kurang variasi dalam kegiatan harian'),
    (uuid_generate_v4(), 'time_management', 'Rekomendasi manajemen waktu', 'Terlalu banyak kegiatan dalam waktu bersamaan atau konflik jadwal')
ON CONFLICT (name) DO NOTHING;

