-- Create work_orders table
CREATE TABLE IF NOT EXISTS work_orders (
    id_order VARCHAR(50) PRIMARY KEY,
    id_pelanggan VARCHAR(50),
    id_petugas VARCHAR(50),
    id_kota VARCHAR(50),
    id_prov VARCHAR(50),
    tgl_order TIMESTAMP,
    tgl_diterima TIMESTAMP,
    tgl_jalan TIMESTAMP,
    tgl_di_lokasi TIMESTAMP,
    tgl_pengecekan TIMESTAMP,
    tgl_selesai_pengecekan TIMESTAMP,
    tgl_pembayaran TIMESTAMP,
    tgl_dikerjakan TIMESTAMP,
    tgl_selesai TIMESTAMP,
    nama_pelanggan VARCHAR(255),
    no_telp VARCHAR(20),
    kwh_meter VARCHAR(50),
    status_mcb VARCHAR(50),
    status_padam VARCHAR(50),
    alamat TEXT,
    detail_alamat TEXT,
    tgl_kunjungan TIMESTAMP,
    pengaduan TEXT,
    catatan TEXT,
    lat_pelanggan DECIMAL(10, 8),
    long_pelanggan DECIMAL(11, 8),
    biaya_transport DECIMAL(12, 2),
    biaya_pemeriksaan DECIMAL(12, 2),
    jumlah_rating INTEGER,
    tgl_rating TIMESTAMP,
    komentar_rating TEXT,
    metode_bayar VARCHAR(50),
    tipe VARCHAR(50),
    status INTEGER,
    status_transfer BOOLEAN DEFAULT false,
    penyebab_batal TEXT,
    keterangan TEXT,
    kd_unit VARCHAR(50),
    kode_pos VARCHAR(10),
    sumber_wo VARCHAR(50),
    last_callback TIMESTAMP,
    ket_callback TEXT,
    is_berulang BOOLEAN DEFAULT false,
    flag_payment BOOLEAN DEFAULT false,
    sent_invoice BOOLEAN DEFAULT false,
    wo_apkt VARCHAR(50),
    source_wo VARCHAR(50),
    email VARCHAR(255),
    nama_provinsi VARCHAR(100),
    ket_callback_dtl TEXT,
    app_version VARCHAR(50),
    autodispatch INTEGER DEFAULT 0,
    ket_dispatch TEXT,
    clear_temper TEXT,
    radius DECIMAL(10, 2),
    response TEXT,
    recovery TEXT,
    sent_invoice_bck1707 BOOLEAN DEFAULT false,
    ba_cancel TEXT,
    flag_bagi_hasil BOOLEAN DEFAULT false,
    flag_potong_deposit BOOLEAN DEFAULT false,
    flag_tambah_dompet BOOLEAN DEFAULT false,
    is_paid_mitra BOOLEAN DEFAULT false,
    callback_by VARCHAR(50),
    kepuasan_pelanggan_respon_layanan TEXT,
    kepuasan_pelanggan_kualitas_petugas TEXT,
    kepuasan_pelanggan_etika_petugas TEXT,
    kepuasan_pelanggan_harga_layanan TEXT,
    kepuasan_pelanggan_temuan TEXT,
    kepuasan_pelanggan_pelanggaran TEXT,
    kepuasan_pelanggan_harga_pelanggan TEXT,
    kepuasan_pelanggan_keterangan TEXT,
    kepuasan_pelanggan_tdk_bersedia TEXT,
    kepuasan_pelanggan_alasan_gagal TEXT,
    kepuasan_pelanggan_ba_klarifikasi TEXT,
    kepuasan_pelanggan_ba_keterangan TEXT,
    kepuasan_pelanggan_status VARCHAR(50),
    dispatch_counter INTEGER DEFAULT 0,
    dispatch_status VARCHAR(50),
    dispatch_log_time TIMESTAMP,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    no_agenda VARCHAR(50),
    is_marketplace BOOLEAN DEFAULT false,
    konfirmasi VARCHAR(10),
    tgl_konfirmasi TIMESTAMP,
    lapor_ulang INTEGER DEFAULT 0,
    waktu_kunjungan TIMESTAMP,
    send_whatsapp BOOLEAN DEFAULT false,
    is_confirmed BOOLEAN DEFAULT false,
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_work_orders_id_pelanggan ON work_orders(id_pelanggan);
CREATE INDEX IF NOT EXISTS idx_work_orders_no_telp ON work_orders(no_telp);
CREATE INDEX IF NOT EXISTS idx_work_orders_tgl_kunjungan ON work_orders(tgl_kunjungan);
CREATE INDEX IF NOT EXISTS idx_work_orders_status ON work_orders(status);
CREATE INDEX IF NOT EXISTS idx_work_orders_created_date ON work_orders(created_date);
CREATE INDEX IF NOT EXISTS idx_work_orders_send_whatsapp ON work_orders(send_whatsapp);

-- Create trigger for updated_at
DROP TRIGGER IF EXISTS update_work_orders_updated_at ON work_orders;
CREATE TRIGGER update_work_orders_updated_at
    BEFORE UPDATE ON work_orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comment
COMMENT ON TABLE work_orders IS 'Tabel untuk menyimpan data work order/service order';
COMMENT ON COLUMN work_orders.id_order IS 'ID unik work order';
COMMENT ON COLUMN work_orders.tgl_kunjungan IS 'Tanggal kunjungan yang dijadwalkan';
COMMENT ON COLUMN work_orders.send_whatsapp IS 'Flag apakah sudah dikirim notifikasi WhatsApp';

