-- Script untuk import work orders dari CSV atau data external
-- Contoh penggunaan dengan COPY command PostgreSQL

-- Jika data dalam format CSV, gunakan:
-- COPY work_orders (
--     id_order, id_pelanggan, id_petugas, id_kota, id_prov,
--     tgl_order, tgl_diterima, tgl_jalan, tgl_di_lokasi, tgl_pengecekan,
--     tgl_selesai_pengecekan, tgl_pembayaran, tgl_dikerjakan, tgl_selesai,
--     nama_pelanggan, no_telp, kwh_meter, status_mcb, status_padam,
--     alamat, detail_alamat, tgl_kunjungan, pengaduan, catatan,
--     lat_pelanggan, long_pelanggan, biaya_transport, biaya_pemeriksaan,
--     jumlah_rating, tgl_rating, komentar_rating, metode_bayar, tipe, status,
--     status_transfer, penyebab_batal, keterangan, kd_unit, kode_pos,
--     sumber_wo, last_callback, ket_callback, is_berulang, flag_payment,
--     sent_invoice, wo_apkt, source_wo, email, nama_provinsi,
--     ket_callback_dtl, app_version, autodispatch, ket_dispatch, clear_temper,
--     radius, response, recovery, sent_invoice_bck1707, ba_cancel,
--     flag_bagi_hasil, flag_potong_deposit, flag_tambah_dompet, is_paid_mitra,
--     callback_by, kepuasan_pelanggan_respon_layanan, kepuasan_pelanggan_kualitas_petugas,
--     kepuasan_pelanggan_etika_petugas, kepuasan_pelanggan_harga_layanan,
--     kepuasan_pelanggan_temuan, kepuasan_pelanggan_pelanggaran,
--     kepuasan_pelanggan_harga_pelanggan, kepuasan_pelanggan_keterangan,
--     kepuasan_pelanggan_tdk_bersedia, kepuasan_pelanggan_alasan_gagal,
--     kepuasan_pelanggan_ba_klarifikasi, kepuasan_pelanggan_ba_keterangan,
--     kepuasan_pelanggan_status, dispatch_counter, dispatch_status, dispatch_log_time,
--     created_date, no_agenda, is_marketplace, konfirmasi, tgl_konfirmasi,
--     lapor_ulang, waktu_kunjungan, send_whatsapp, is_confirmed, created_by
-- )
-- FROM '/path/to/work_orders.csv'
-- WITH (FORMAT csv, HEADER true, DELIMITER '|');

-- Contoh insert manual:
INSERT INTO work_orders (
    id_order, id_pelanggan, nama_pelanggan, no_telp,
    tgl_order, tgl_kunjungan, pengaduan, catatan,
    alamat, detail_alamat, lat_pelanggan, long_pelanggan,
    kwh_meter, status_mcb, status_padam,
    biaya_transport, biaya_pemeriksaan,
    metode_bayar, tipe, status,
    nama_provinsi, kd_unit, kode_pos,
    sumber_wo, source_wo, email,
    send_whatsapp, is_confirmed, created_by, created_date
) VALUES (
    '30292',
    '542104357017',
    'yenny kumala',
    '62818992459',
    '2025-12-16 06:37:17.667',
    '2025-12-16 08:36:00.000',
    'MCB turun pada saat charging',
    'MCB turun pada saat charging',
    'JL V-PLUIT BARAT 2 NO 4 , PLUIT, PENJARINGA No.',
    'JL V-PLUIT BARAT 2 NO 4 , PLUIT, PENJARINGA No.',
    -6.1251113386929745,
    106.78731624037027,
    'Pra Bayar',
    'on',
    'ya',
    20000,
    50000,
    NULL,
    'PLN MOBILE',
    -2,
    'DKI Jakarta',
    '14320',
    NULL,
    NULL,
    'PLN MOBILE',
    'yenny_in2@yahoo.com',
    false,
    false,
    'OLD Service',
    '2025-12-16 06:37:17.667'
)
ON CONFLICT (id_order) DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP;

