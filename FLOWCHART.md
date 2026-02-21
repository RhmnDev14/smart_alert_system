# Flowchart Smart Alert System

## 1. Flowchart Utama Sistem

```mermaid
flowchart TD
    Start([Sistem Dimulai]) --> Init[Inisialisasi Server & Telegram Bot]
    Init --> Scheduler[Setup Scheduler Dinamis]

    Scheduler --> RealtimeReminder[Activity Reminder - Cek per 1 Menit]
    Scheduler --> MorningAlert[Alert Pagi - Ambil Konfig dari .env]
    Scheduler --> EveningSummary[Summary Malam - Ambil Konfig dari .env]

    Init --> TelegramListener[Listener Telegram Messages via Long-polling]

    TelegramListener --> CheckUser{User Baru?}
    CheckUser -->|Ya| SendWelcome[Kirim Pesan Welcome/Default]
    CheckUser -->|Tidak| AIGateway[Kirim Pesan ke AI Gateway]

    SendWelcome --> EndMsg([Selesai])

    AIGateway --> ProcessMsg[Proses Obrolan Natural & Analisis Jadwal]
    ProcessMsg --> CheckSchedule{Ada Konteks Jadwal?}

    CheckSchedule -->|Ya| ExtractData[Mengekstrak JSON Jadwal Event]
    CheckSchedule -->|Tidak| SkipExtract[Skipping - Hanya Murni Obrolan]

    ExtractData --> InsertDB[Simpan Jadwal ke DB secara Sinkron]
    InsertDB --> SendAIResponse[Telegram Kirim Respons AI ke User]
    SkipExtract --> SendAIResponse
    SendAIResponse --> EndMsg

    MorningAlert --> ProcessMorning[AI Generate Ucapan Motivasi, Rekap Jadwal & Menanya Rencana Tambahan]
    ProcessMorning --> SendMorningAlert[Telegram Kirim Alert Pagi ke User]
    SendMorningAlert --> EndAlert([Selesai])

    EveningSummary --> ProcessEvening[AI Generate Ringkasan Selesai, Analisa Pola & Insight Kesehatan]
    ProcessEvening --> SendEveningSummary[Telegram Kirim Summary Malam ke User]
    SendEveningSummary --> EndSummary([Selesai])

    RealtimeReminder --> FetchPending[Ambil Aktivitas Pending yang Tiba Masanya]
    FetchPending --> SendReminder[Kirim Notifikasi via Telegram]
    SendReminder --> EndReminder([Selesai])
```

## 2. Flowchart Proses Gateway AI

```mermaid
flowchart TD
    Start([Pesan Masuk]) --> InsertSystemPrompt[Sisipkan Datetime & Instruksi Sistem]
    InsertSystemPrompt --> AIProcessing[Model AI Menganalisa & Merespons Pesan]

    AIProcessing --> OutputFormat{Return Valid JSON?}

    OutputFormat -->|Tidak Valid| FallbackResponse[Ambil Raw String sbg Response Obrolan Biasa]
    OutputFormat -->|Valid| Extract[Ekstrak obj JSON: Response & Schedule]

    Extract --> HasSchedule{has_schedule == true?}
    HasSchedule -->|Tidak| OnlyResponse[Kembalikan Reponse Saja]
    HasSchedule -->|Ya| ExtractScheduleData[Ekstrak: Judul, Deskripsi, Schedule_Time, Priority]

    OnlyResponse --> End([Selesai - Gateway Mengembalikan Obrolan])
    ExtractScheduleData --> End
    FallbackResponse --> End
```

## 3. Flowchart Scheduler Alert Pagi

```mermaid
flowchart TD
    Start([Jam Konfigurasi .env]) --> GetAllUsers[Ambil Semua User Aktif]
    GetAllUsers --> LoopUser{Ada User?}

    LoopUser -->|Ya| GetUserActivities[Ambil Kegiatan User Hari Ini]
    LoopUser -->|Tidak| End([Selesai])

    GetUserActivities --> GetHealthContext[Ambil Konteks Kesehatan User]

    GetHealthContext --> AnalyzeWithAI[Analisis dengan AI]
    AnalyzeWithAI --> GenerateAlert[Generate Ucapan Motivasi, Rekap Acara, Nanya Rencana Tambahan]

    GenerateAlert --> SendAlert[Kirim Alert via Telegram ke User]

    SendAlert --> NextUser[User Berikutnya]
    NextUser --> LoopUser
```

## 4. Flowchart Scheduler Summary Malam

```mermaid
flowchart TD
    Start([Jam Konfigurasi .env]) --> GetAllUsers[Ambil Semua User Aktif]
    GetAllUsers --> LoopUser{Ada User?}

    LoopUser -->|Ya| GetTodayActivities[Ambil Kegiatan Harian Status Completed]
    LoopUser -->|Tidak| End([Selesai])

    GetTodayActivities --> AIAnalysis[Disalurkan ke AI]
    AIAnalysis --> GenerateInsights[AI Generate Obrolan, Ringkasan, Pola Produktivitas, Feedback & Rekomendasi Hari Esok]

    GenerateInsights --> SendSummary[Kirim Summary via Telegram ke User]

    SendSummary --> NextUser[User Berikutnya]
    NextUser --> LoopUser
```

## 5. Flowchart Pengingat Real-Time (Per-Menit)

```mermaid
flowchart TD
    Start([Cron: Setiap Menit]) --> GetPending[Ambil Aktivitas Status 'Pending' dg Time <= Now()]
    GetPending --> LoopActivity{Ada Kegiatan Expired?}

    LoopActivity -->|Tidak| End([Selesai])
    LoopActivity -->|Ya| CheckReminded{Sudah Diremind?}

    CheckReminded -->|Ya| NextActivity
    CheckReminded -->|Tidak| GetUserDetail[Ambil Profil User Terkait]

    GetUserDetail --> FormatMsg[Format Notifikasi: Pengingat, Judul, Waktu, Catatan]
    FormatMsg --> SendTelegram[Kirim Notifikasi via Telegram]

    SendTelegram --> UpdateStatus[Update Kolom `ReminderTime`]
    UpdateStatus --> NextActivity[Kegiatan Berikutnya]

    NextActivity --> LoopActivity
```

## 6. Flowchart Sistem AI untuk Rekomendasi Kesehatan

```mermaid
flowchart TD
    Start([Trigger AI Health]) --> GetUserProfile[Ambil Profil User]
    GetUserProfile --> GetActivityHistory[Ambil History Kegiatan]
    GetActivityHistory --> GetHealthData[Ambil Data Kesehatan User]

    GetHealthData --> AnalyzePattern[Analisis Pola Kegiatan]
    AnalyzePattern --> IdentifyHealthIssues[Identifikasi Masalah Kesehatan Potensial]

    IdentifyHealthIssues --> Contextualize[Kontekstualisasi dengan Kegiatan]
    Contextualize --> GenerateRecommendations[Generate Rekomendasi Spesifik]

    GenerateRecommendations --> FormatRecommendation[Format Rekomendasi]
    FormatRecommendation --> ReturnRecommendation[Return Rekomendasi]
    ReturnRecommendation --> End([Selesai])
```
