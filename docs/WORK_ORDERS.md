# Work Orders Documentation

## Overview

Fitur Work Orders memungkinkan sistem untuk menyimpan dan mengelola data work order/service order dari sistem eksternal (misalnya sistem PLN atau layanan listrik).

## Database Schema

Tabel `work_orders` memiliki 90+ kolom yang mencakup:
- Informasi dasar: `id_order`, `id_pelanggan`, `nama_pelanggan`, `no_telp`
- Tanggal-tanggal penting: `tgl_order`, `tgl_kunjungan`, `tgl_selesai`, dll
- Lokasi: `alamat`, `lat_pelanggan`, `long_pelanggan`
- Status dan tracking: `status`, `send_whatsapp`, `is_confirmed`
- Biaya: `biaya_transport`, `biaya_pemeriksaan`
- Rating dan kepuasan pelanggan
- Dan banyak lagi...

## Migration

Jalankan migration untuk membuat tabel:

```bash
make migrate
```

Atau secara manual:

```bash
./bin/migrate
```

File migration: `migrations/020_create_work_orders_table.sql`

## Entity

Entity `WorkOrder` berada di `internal/domain/entity/work_order.go` dengan helper methods:
- `GetTglKunjungan()`: Mengembalikan tanggal kunjungan (prioritas `tgl_kunjungan`, lalu `waktu_kunjungan`)
- `GetNoTelpClean()`: Membersihkan nomor telepon (menghapus prefix country code)
- `IsPending()`: Mengecek apakah work order masih pending
- `IsCompleted()`: Mengecek apakah work order sudah selesai

## Repository

Repository interface: `internal/domain/repository/work_order_repository.go`

Implementasi: `internal/infrastructure/repository/work_order_repository_impl.go`

### Methods yang tersedia:

- `Create(ctx, workOrder)`: Membuat work order baru
- `GetByID(ctx, idOrder)`: Mengambil work order berdasarkan ID
- `GetByNoTelp(ctx, noTelp)`: Mengambil semua work order berdasarkan nomor telepon
- `GetByTglKunjungan(ctx, date)`: Mengambil work order yang dijadwalkan pada tanggal tertentu
- `GetPending(ctx)`: Mengambil semua work order yang masih pending
- `GetPendingByNoTelp(ctx, noTelp)`: Mengambil work order pending untuk nomor telepon tertentu
- `Update(ctx, workOrder)`: Update work order
- `UpdateSendWhatsapp(ctx, idOrder, sent)`: Update flag `send_whatsapp`
- `Delete(ctx, idOrder)`: Hapus work order
- `BulkCreate(ctx, workOrders)`: Import multiple work orders sekaligus

## Import Data

### 1. Menggunakan Go Tool

Build tool import:

```bash
go build -o bin/import_work_orders cmd/import_work_orders/main.go
```

Jalankan import dari file CSV (format pipe-delimited):

```bash
./bin/import_work_orders work_orders.csv
```

Format CSV harus menggunakan pipe (`|`) sebagai delimiter dan memiliki header yang sesuai dengan kolom database.

### 2. Menggunakan SQL Script

Lihat contoh di `scripts/import_work_orders.sql` untuk:
- Menggunakan PostgreSQL `COPY` command
- Insert manual dengan `INSERT INTO`

## Integrasi dengan Smart Alert System

Work orders dapat diintegrasikan dengan Smart Alert System untuk:

1. **Mengkonversi Work Order menjadi Activity**
   - Gunakan `tgl_kunjungan` sebagai `scheduled_time`
   - Gunakan `pengaduan` atau `catatan` sebagai `description`
   - Link dengan user berdasarkan `no_telp`

2. **Notifikasi WhatsApp**
   - Kirim reminder sebelum `tgl_kunjungan`
   - Update flag `send_whatsapp` setelah mengirim notifikasi

3. **Tracking Status**
   - Monitor work order yang pending
   - Update status setelah selesai

## Contoh Penggunaan

### Mengambil Work Order Pending untuk Nomor Telepon

```go
workOrderRepo := infraRepo.NewWorkOrderRepository(db)
pendingOrders, err := workOrderRepo.GetPendingByNoTelp(ctx, "62818992459")
```

### Mengkonversi Work Order menjadi Activity

```go
workOrder, _ := workOrderRepo.GetByID(ctx, "30292")
if workOrder != nil && workOrder.GetTglKunjungan() != nil {
    activity := &entity.Activity{
        UserID: userID,
        Title: workOrder.Pengaduan.String,
        Description: workOrder.Catatan.String,
        ScheduledTime: *workOrder.GetTglKunjungan(),
        Status: "pending",
    }
    // Save activity...
}
```

## Indexes

Tabel memiliki indexes untuk performa query:
- `idx_work_orders_id_pelanggan`: Query berdasarkan ID pelanggan
- `idx_work_orders_no_telp`: Query berdasarkan nomor telepon
- `idx_work_orders_tgl_kunjungan`: Query berdasarkan tanggal kunjungan
- `idx_work_orders_status`: Query berdasarkan status
- `idx_work_orders_created_date`: Query berdasarkan tanggal dibuat
- `idx_work_orders_send_whatsapp`: Query work order yang belum dikirim notifikasi

## Catatan Penting

1. **Format Nomor Telepon**: Sistem menggunakan format dengan country code (misalnya `62818992459`). Gunakan `GetNoTelpClean()` untuk mendapatkan format tanpa country code.

2. **Status**: 
   - Status negatif (< 0) = Pending
   - Status positif (> 0) dengan `tgl_selesai` = Completed

3. **Tanggal Kunjungan**: Sistem memprioritaskan `tgl_kunjungan`, jika tidak ada maka menggunakan `waktu_kunjungan`.

4. **Idempotent**: Method `Create` menggunakan `ON CONFLICT DO UPDATE`, sehingga import ulang tidak akan membuat duplikat.

## Next Steps

1. Buat use case untuk mengkonversi work order menjadi activity
2. Integrasikan dengan scheduler untuk mengirim reminder
3. Buat handler untuk menerima webhook dari sistem eksternal
4. Buat API endpoint untuk query dan update work orders

