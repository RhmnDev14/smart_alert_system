# Flowchart Smart Alert System

## 1. Flowchart Utama Sistem

```mermaid
flowchart TD
    Start([Sistem Dimulai]) --> Init[Inisialisasi Server Waha]
    Init --> Scheduler[Setup Scheduler]
    Scheduler --> MorningAlert[Alert Pagi Jam 05:00]
    Scheduler --> EveningSummary[Summary Malam Jam 22:00]
    
    Init --> WhatsAppListener[Listener WhatsApp Messages]
    
    WhatsAppListener --> CheckUser{User Baru?}
    CheckUser -->|Ya| SendWelcome[Kirim Pesan Welcome/Default]
    CheckUser -->|Tidak| ParseMessage[Parse Pesan User]
    
    SendWelcome --> EndMsg([Selesai])
    
    ParseMessage --> AIProcess[Proses dengan AI]
    AIProcess --> ExtractIntent[Ekstrak Intent & Data]
    ExtractIntent --> IntentType{Type Intent?}
    
    IntentType -->|Tambah Kegiatan| SaveActivity[Simpan Kegiatan ke DB]
    IntentType -->|Hapus Kegiatan| DeleteActivity[Hapus Kegiatan dari DB]
    IntentType -->|Lihat Kegiatan| GetActivities[Ambil Daftar Kegiatan]
    IntentType -->|Update Kegiatan| UpdateActivity[Update Kegiatan]
    IntentType -->|Pertanyaan Umum| GenerateResponse[Generate Response AI]
    
    SaveActivity --> ConfirmSave[Konfirmasi Kegiatan Tersimpan]
    DeleteActivity --> ConfirmDelete[Konfirmasi Kegiatan Dihapus]
    GetActivities --> FormatActivities[Format Daftar Kegiatan]
    UpdateActivity --> ConfirmUpdate[Konfirmasi Kegiatan Diupdate]
    GenerateResponse --> SendResponse[Kirim Response]
    
    ConfirmSave --> EndMsg
    ConfirmDelete --> EndMsg
    FormatActivities --> EndMsg
    ConfirmUpdate --> EndMsg
    SendResponse --> EndMsg
    
    MorningAlert --> GetTodayActivities[Ambil Kegiatan Hari Ini]
    GetTodayActivities --> GenerateHealthTips[Generate Tips Kesehatan AI]
    GenerateHealthTips --> FormatMorningAlert[Format Alert Pagi]
    FormatMorningAlert --> SendMorningAlert[Kirim Alert ke User]
    SendMorningAlert --> EndAlert([Selesai])
    
    EveningSummary --> GetTodayCompleted[Ambil Kegiatan Selesai Hari Ini]
    GetTodayCompleted --> AnalyzeActivities[Analisis Kegiatan dengan AI]
    AnalyzeActivities --> GenerateSummary[Generate Summary & Rekomendasi]
    GenerateSummary --> FormatEveningSummary[Format Summary Malam]
    FormatEveningSummary --> SendEveningSummary[Kirim Summary ke User]
    SendEveningSummary --> EndSummary([Selesai])
```

## 2. Flowchart Proses AI untuk Parsing Pesan

```mermaid
flowchart TD
    Start([Pesan Masuk]) --> Preprocess[Preprocess Pesan]
    Preprocess --> AIExtract[AI Extract Intent & Entities]
    AIExtract --> Intent{Intent Terdeteksi?}
    
    Intent -->|Tambah Kegiatan| ExtractActivity[Extract: Nama, Waktu, Deskripsi]
    Intent -->|Hapus Kegiatan| ExtractID[Extract: ID Kegiatan]
    Intent -->|Update Kegiatan| ExtractUpdate[Extract: ID, Field, Value]
    Intent -->|Lihat Kegiatan| ExtractDate[Extract: Tanggal Optional]
    Intent -->|Pertanyaan| ExtractQuestion[Extract: Pertanyaan]
    Intent -->|Tidak Jelas| Fallback[Fallback: Tanya Kembali]
    
    ExtractActivity --> ValidateActivity{Valid?}
    ExtractID --> ValidateID{Valid?}
    ExtractUpdate --> ValidateUpdate{Valid?}
    ExtractDate --> GetActivities[Ambil Kegiatan]
    ExtractQuestion --> ProcessQuestion[Proses Pertanyaan]
    Fallback --> SendClarification[Kirim Pesan Klarifikasi]
    
    ValidateActivity -->|Ya| ReturnActivity[Return Activity Data]
    ValidateActivity -->|Tidak| SendError[Kirim Error Message]
    
    ValidateID -->|Ya| ReturnID[Return ID]
    ValidateID -->|Tidak| SendError
    
    ValidateUpdate -->|Ya| ReturnUpdate[Return Update Data]
    ValidateUpdate -->|Tidak| SendError
    
    GetActivities --> ReturnActivities[Return Activities]
    ProcessQuestion --> ReturnAnswer[Return Answer]
    
    ReturnActivity --> End([Selesai])
    ReturnID --> End
    ReturnUpdate --> End
    ReturnActivities --> End
    ReturnAnswer --> End
    SendError --> End
    SendClarification --> End
```

## 3. Flowchart Scheduler Alert Pagi

```mermaid
flowchart TD
    Start([Jam 05:00]) --> GetAllUsers[Ambil Semua User Aktif]
    GetAllUsers --> LoopUser{Ada User?}
    
    LoopUser -->|Ya| GetUserActivities[Ambil Kegiatan User Hari Ini]
    LoopUser -->|Tidak| End([Selesai])
    
    GetUserActivities --> CheckActivities{Ada Kegiatan?}
    
    CheckActivities -->|Ya| GetHealthContext[Ambil Konteks Kesehatan User]
    CheckActivities -->|Tidak| GenerateGeneralTips[Generate Tips Umum]
    
    GetHealthContext --> AnalyzeWithAI[Analisis dengan AI]
    AnalyzeWithAI --> GeneratePersonalizedTips[Generate Tips Personalisasi]
    GenerateGeneralTips --> FormatGeneralAlert[Format Alert Umum]
    GeneratePersonalizedTips --> FormatPersonalizedAlert[Format Alert Personalisasi]
    
    FormatGeneralAlert --> SendAlert[Kirim Alert ke User]
    FormatPersonalizedAlert --> SendAlert
    
    SendAlert --> NextUser[User Berikutnya]
    NextUser --> LoopUser
```

## 4. Flowchart Scheduler Summary Malam

```mermaid
flowchart TD
    Start([Jam 22:00]) --> GetAllUsers[Ambil Semua User Aktif]
    GetAllUsers --> LoopUser{Ada User?}
    
    LoopUser -->|Ya| GetTodayActivities[Ambil Kegiatan Hari Ini]
    LoopUser -->|Tidak| End([Selesai])
    
    GetTodayActivities --> CategorizeActivities[Kategorisasi Kegiatan]
    CategorizeActivities --> AnalyzeCompletion[Analisis Tingkat Penyelesaian]
    AnalyzeCompletion --> GetHealthPatterns[Analisis Pola Kesehatan]
    
    GetHealthPatterns --> AIAnalysis[Analisis dengan AI]
    AIAnalysis --> GenerateInsights[Generate Insights]
    GenerateInsights --> GenerateRecommendations[Generate Rekomendasi]
    
    GenerateRecommendations --> FormatSummary[Format Summary]
    FormatSummary --> SendSummary[Kirim Summary ke User]
    
    SendSummary --> NextUser[User Berikutnya]
    NextUser --> LoopUser
```

## 5. Flowchart Proses Input Kegiatan User

```mermaid
flowchart TD
    Start([User Input Kegiatan]) --> ReceiveMessage[Terima Pesan WhatsApp]
    ReceiveMessage --> AIParse[AI Parse Pesan]
    
    AIParse --> ExtractInfo[Extract: Nama, Waktu, Deskripsi, Kategori]
    ExtractInfo --> ValidateTime{Waktu Valid?}
    
    ValidateTime -->|Tidak| AskTime[Kirim Pesan: Mohon Spesifik Waktu]
    AskTime --> WaitResponse[Tunggu Response User]
    WaitResponse --> ReceiveMessage
    
    ValidateTime -->|Ya| CheckConflict{Cek Konflik Waktu?}
    CheckConflict -->|Ada Konflik| NotifyConflict[Notifikasi Konflik]
    NotifyConflict --> AskConfirm{Tanya Konfirmasi}
    
    AskConfirm -->|Ya| SaveActivity[Simpan Kegiatan]
    AskConfirm -->|Tidak| CancelSave[Batal Simpan]
    
    CheckConflict -->|Tidak Ada| SaveActivity
    
    SaveActivity --> GenerateConfirmation[Generate Konfirmasi]
    GenerateConfirmation --> SendConfirmation[Kirim Konfirmasi ke User]
    
    CancelSave --> End([Selesai])
    SendConfirmation --> End
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

