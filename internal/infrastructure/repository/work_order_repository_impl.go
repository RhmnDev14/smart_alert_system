package repository

import (
	"context"
	"database/sql"
	"fmt"

	"smart_alert_system/internal/domain/entity"
)

type WorkOrderRepositoryImpl struct {
	db *sql.DB
}

func NewWorkOrderRepository(db *sql.DB) *WorkOrderRepositoryImpl {
	return &WorkOrderRepositoryImpl{db: db}
}

func (r *WorkOrderRepositoryImpl) Create(ctx context.Context, workOrder *entity.WorkOrder) error {
	query := `
		INSERT INTO work_orders (
			id_order, id_pelanggan, id_petugas, id_kota, id_prov,
			tgl_order, tgl_diterima, tgl_jalan, tgl_di_lokasi, tgl_pengecekan,
			tgl_selesai_pengecekan, tgl_pembayaran, tgl_dikerjakan, tgl_selesai,
			nama_pelanggan, no_telp, kwh_meter, status_mcb, status_padam,
			alamat, detail_alamat, tgl_kunjungan, pengaduan, catatan,
			lat_pelanggan, long_pelanggan, biaya_transport, biaya_pemeriksaan,
			jumlah_rating, tgl_rating, komentar_rating, metode_bayar, tipe, status,
			status_transfer, penyebab_batal, keterangan, kd_unit, kode_pos,
			sumber_wo, last_callback, ket_callback, is_berulang, flag_payment,
			sent_invoice, wo_apkt, source_wo, email, nama_provinsi,
			ket_callback_dtl, app_version, autodispatch, ket_dispatch, clear_temper,
			radius, response, recovery, sent_invoice_bck1707, ba_cancel,
			flag_bagi_hasil, flag_potong_deposit, flag_tambah_dompet, is_paid_mitra,
			callback_by, kepuasan_pelanggan_respon_layanan, kepuasan_pelanggan_kualitas_petugas,
			kepuasan_pelanggan_etika_petugas, kepuasan_pelanggan_harga_layanan,
			kepuasan_pelanggan_temuan, kepuasan_pelanggan_pelanggaran,
			kepuasan_pelanggan_harga_pelanggan, kepuasan_pelanggan_keterangan,
			kepuasan_pelanggan_tdk_bersedia, kepuasan_pelanggan_alasan_gagal,
			kepuasan_pelanggan_ba_klarifikasi, kepuasan_pelanggan_ba_keterangan,
			kepuasan_pelanggan_status, dispatch_counter, dispatch_status, dispatch_log_time,
			created_date, no_agenda, is_marketplace, konfirmasi, tgl_konfirmasi,
			lapor_ulang, waktu_kunjungan, send_whatsapp, is_confirmed, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40,
			$41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58, $59, $60,
			$61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76, $77, $78, $79, $80,
			$81, $82, $83, $84, $85, $86, $87, $88, $89, $90
		)
		ON CONFLICT (id_order) DO UPDATE SET
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.ExecContext(ctx, query,
		workOrder.IDOrder, workOrder.IDPelanggan, workOrder.IDPetugas, workOrder.IDKota, workOrder.IDProv,
		workOrder.TglOrder, workOrder.TglDiterima, workOrder.TglJalan, workOrder.TglDiLokasi, workOrder.TglPengecekan,
		workOrder.TglSelesaiPengecekan, workOrder.TglPembayaran, workOrder.TglDikerjakan, workOrder.TglSelesai,
		workOrder.NamaPelanggan, workOrder.NoTelp, workOrder.KwhMeter, workOrder.StatusMcb, workOrder.StatusPadam,
		workOrder.Alamat, workOrder.DetailAlamat, workOrder.TglKunjungan, workOrder.Pengaduan, workOrder.Catatan,
		workOrder.LatPelanggan, workOrder.LongPelanggan, workOrder.BiayaTransport, workOrder.BiayaPemeriksaan,
		workOrder.JumlahRating, workOrder.TglRating, workOrder.KomentarRating, workOrder.MetodeBayar, workOrder.Tipe, workOrder.Status,
		workOrder.StatusTransfer, workOrder.PenyebabBatal, workOrder.Keterangan, workOrder.KdUnit, workOrder.KodePos,
		workOrder.SumberWo, workOrder.LastCallback, workOrder.KetCallback, workOrder.IsBerulang, workOrder.FlagPayment,
		workOrder.SentInvoice, workOrder.WoApkt, workOrder.SourceWo, workOrder.Email, workOrder.NamaProvinsi,
		workOrder.KetCallbackDtl, workOrder.AppVersion, workOrder.Autodispatch, workOrder.KetDispatch, workOrder.ClearTemper,
		workOrder.Radius, workOrder.Response, workOrder.Recovery, workOrder.SentInvoiceBck1707, workOrder.BaCancel,
		workOrder.FlagBagiHasil, workOrder.FlagPotongDeposit, workOrder.FlagTambahDompet, workOrder.IsPaidMitra,
		workOrder.CallbackBy, workOrder.KepuasanPelangganResponLayanan, workOrder.KepuasanPelangganKualitasPetugas,
		workOrder.KepuasanPelangganEtikaPetugas, workOrder.KepuasanPelangganHargaLayanan,
		workOrder.KepuasanPelangganTemuan, workOrder.KepuasanPelangganPelanggaran,
		workOrder.KepuasanPelangganHargaPelanggan, workOrder.KepuasanPelangganKeterangan,
		workOrder.KepuasanPelangganTdkBersedia, workOrder.KepuasanPelangganAlasanGagal,
		workOrder.KepuasanPelangganBaKlarifikasi, workOrder.KepuasanPelangganBaKeterangan,
		workOrder.KepuasanPelangganStatus, workOrder.DispatchCounter, workOrder.DispatchStatus, workOrder.DispatchLogTime,
		workOrder.CreatedDate, workOrder.NoAgenda, workOrder.IsMarketplace, workOrder.Konfirmasi, workOrder.TglKonfirmasi,
		workOrder.LaporUlang, workOrder.WaktuKunjungan, workOrder.SendWhatsapp, workOrder.IsConfirmed, workOrder.CreatedBy,
	)

	return err
}

func (r *WorkOrderRepositoryImpl) GetByID(ctx context.Context, idOrder string) (*entity.WorkOrder, error) {
	query := `SELECT * FROM work_orders WHERE id_order = $1`
	wo := &entity.WorkOrder{}
	err := r.db.QueryRowContext(ctx, query, idOrder).Scan(
		&wo.IDOrder, &wo.IDPelanggan, &wo.IDPetugas, &wo.IDKota, &wo.IDProv,
		&wo.TglOrder, &wo.TglDiterima, &wo.TglJalan, &wo.TglDiLokasi, &wo.TglPengecekan,
		&wo.TglSelesaiPengecekan, &wo.TglPembayaran, &wo.TglDikerjakan, &wo.TglSelesai,
		&wo.NamaPelanggan, &wo.NoTelp, &wo.KwhMeter, &wo.StatusMcb, &wo.StatusPadam,
		&wo.Alamat, &wo.DetailAlamat, &wo.TglKunjungan, &wo.Pengaduan, &wo.Catatan,
		&wo.LatPelanggan, &wo.LongPelanggan, &wo.BiayaTransport, &wo.BiayaPemeriksaan,
		&wo.JumlahRating, &wo.TglRating, &wo.KomentarRating, &wo.MetodeBayar, &wo.Tipe, &wo.Status,
		&wo.StatusTransfer, &wo.PenyebabBatal, &wo.Keterangan, &wo.KdUnit, &wo.KodePos,
		&wo.SumberWo, &wo.LastCallback, &wo.KetCallback, &wo.IsBerulang, &wo.FlagPayment,
		&wo.SentInvoice, &wo.WoApkt, &wo.SourceWo, &wo.Email, &wo.NamaProvinsi,
		&wo.KetCallbackDtl, &wo.AppVersion, &wo.Autodispatch, &wo.KetDispatch, &wo.ClearTemper,
		&wo.Radius, &wo.Response, &wo.Recovery, &wo.SentInvoiceBck1707, &wo.BaCancel,
		&wo.FlagBagiHasil, &wo.FlagPotongDeposit, &wo.FlagTambahDompet, &wo.IsPaidMitra,
		&wo.CallbackBy, &wo.KepuasanPelangganResponLayanan, &wo.KepuasanPelangganKualitasPetugas,
		&wo.KepuasanPelangganEtikaPetugas, &wo.KepuasanPelangganHargaLayanan,
		&wo.KepuasanPelangganTemuan, &wo.KepuasanPelangganPelanggaran,
		&wo.KepuasanPelangganHargaPelanggan, &wo.KepuasanPelangganKeterangan,
		&wo.KepuasanPelangganTdkBersedia, &wo.KepuasanPelangganAlasanGagal,
		&wo.KepuasanPelangganBaKlarifikasi, &wo.KepuasanPelangganBaKeterangan,
		&wo.KepuasanPelangganStatus, &wo.DispatchCounter, &wo.DispatchStatus, &wo.DispatchLogTime,
		&wo.CreatedDate, &wo.NoAgenda, &wo.IsMarketplace, &wo.Konfirmasi, &wo.TglKonfirmasi,
		&wo.LaporUlang, &wo.WaktuKunjungan, &wo.SendWhatsapp, &wo.IsConfirmed, &wo.CreatedBy,
		&wo.CreatedAt, &wo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return wo, nil
}

func (r *WorkOrderRepositoryImpl) GetByNoTelp(ctx context.Context, noTelp string) ([]*entity.WorkOrder, error) {
	query := `SELECT * FROM work_orders WHERE no_telp = $1 ORDER BY tgl_kunjungan DESC, created_date DESC`
	rows, err := r.db.QueryContext(ctx, query, noTelp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workOrders []*entity.WorkOrder
	for rows.Next() {
		wo := &entity.WorkOrder{}
		err := rows.Scan(
			&wo.IDOrder, &wo.IDPelanggan, &wo.IDPetugas, &wo.IDKota, &wo.IDProv,
			&wo.TglOrder, &wo.TglDiterima, &wo.TglJalan, &wo.TglDiLokasi, &wo.TglPengecekan,
			&wo.TglSelesaiPengecekan, &wo.TglPembayaran, &wo.TglDikerjakan, &wo.TglSelesai,
			&wo.NamaPelanggan, &wo.NoTelp, &wo.KwhMeter, &wo.StatusMcb, &wo.StatusPadam,
			&wo.Alamat, &wo.DetailAlamat, &wo.TglKunjungan, &wo.Pengaduan, &wo.Catatan,
			&wo.LatPelanggan, &wo.LongPelanggan, &wo.BiayaTransport, &wo.BiayaPemeriksaan,
			&wo.JumlahRating, &wo.TglRating, &wo.KomentarRating, &wo.MetodeBayar, &wo.Tipe, &wo.Status,
			&wo.StatusTransfer, &wo.PenyebabBatal, &wo.Keterangan, &wo.KdUnit, &wo.KodePos,
			&wo.SumberWo, &wo.LastCallback, &wo.KetCallback, &wo.IsBerulang, &wo.FlagPayment,
			&wo.SentInvoice, &wo.WoApkt, &wo.SourceWo, &wo.Email, &wo.NamaProvinsi,
			&wo.KetCallbackDtl, &wo.AppVersion, &wo.Autodispatch, &wo.KetDispatch, &wo.ClearTemper,
			&wo.Radius, &wo.Response, &wo.Recovery, &wo.SentInvoiceBck1707, &wo.BaCancel,
			&wo.FlagBagiHasil, &wo.FlagPotongDeposit, &wo.FlagTambahDompet, &wo.IsPaidMitra,
			&wo.CallbackBy, &wo.KepuasanPelangganResponLayanan, &wo.KepuasanPelangganKualitasPetugas,
			&wo.KepuasanPelangganEtikaPetugas, &wo.KepuasanPelangganHargaLayanan,
			&wo.KepuasanPelangganTemuan, &wo.KepuasanPelangganPelanggaran,
			&wo.KepuasanPelangganHargaPelanggan, &wo.KepuasanPelangganKeterangan,
			&wo.KepuasanPelangganTdkBersedia, &wo.KepuasanPelangganAlasanGagal,
			&wo.KepuasanPelangganBaKlarifikasi, &wo.KepuasanPelangganBaKeterangan,
			&wo.KepuasanPelangganStatus, &wo.DispatchCounter, &wo.DispatchStatus, &wo.DispatchLogTime,
			&wo.CreatedDate, &wo.NoAgenda, &wo.IsMarketplace, &wo.Konfirmasi, &wo.TglKonfirmasi,
			&wo.LaporUlang, &wo.WaktuKunjungan, &wo.SendWhatsapp, &wo.IsConfirmed, &wo.CreatedBy,
			&wo.CreatedAt, &wo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		workOrders = append(workOrders, wo)
	}
	return workOrders, rows.Err()
}

func (r *WorkOrderRepositoryImpl) GetByTglKunjungan(ctx context.Context, date string) ([]*entity.WorkOrder, error) {
	query := `
		SELECT * FROM work_orders 
		WHERE DATE(tgl_kunjungan) = $1 OR DATE(waktu_kunjungan) = $1
		ORDER BY tgl_kunjungan ASC, waktu_kunjungan ASC
	`
	rows, err := r.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workOrders []*entity.WorkOrder
	for rows.Next() {
		wo := &entity.WorkOrder{}
		err := rows.Scan(
			&wo.IDOrder, &wo.IDPelanggan, &wo.IDPetugas, &wo.IDKota, &wo.IDProv,
			&wo.TglOrder, &wo.TglDiterima, &wo.TglJalan, &wo.TglDiLokasi, &wo.TglPengecekan,
			&wo.TglSelesaiPengecekan, &wo.TglPembayaran, &wo.TglDikerjakan, &wo.TglSelesai,
			&wo.NamaPelanggan, &wo.NoTelp, &wo.KwhMeter, &wo.StatusMcb, &wo.StatusPadam,
			&wo.Alamat, &wo.DetailAlamat, &wo.TglKunjungan, &wo.Pengaduan, &wo.Catatan,
			&wo.LatPelanggan, &wo.LongPelanggan, &wo.BiayaTransport, &wo.BiayaPemeriksaan,
			&wo.JumlahRating, &wo.TglRating, &wo.KomentarRating, &wo.MetodeBayar, &wo.Tipe, &wo.Status,
			&wo.StatusTransfer, &wo.PenyebabBatal, &wo.Keterangan, &wo.KdUnit, &wo.KodePos,
			&wo.SumberWo, &wo.LastCallback, &wo.KetCallback, &wo.IsBerulang, &wo.FlagPayment,
			&wo.SentInvoice, &wo.WoApkt, &wo.SourceWo, &wo.Email, &wo.NamaProvinsi,
			&wo.KetCallbackDtl, &wo.AppVersion, &wo.Autodispatch, &wo.KetDispatch, &wo.ClearTemper,
			&wo.Radius, &wo.Response, &wo.Recovery, &wo.SentInvoiceBck1707, &wo.BaCancel,
			&wo.FlagBagiHasil, &wo.FlagPotongDeposit, &wo.FlagTambahDompet, &wo.IsPaidMitra,
			&wo.CallbackBy, &wo.KepuasanPelangganResponLayanan, &wo.KepuasanPelangganKualitasPetugas,
			&wo.KepuasanPelangganEtikaPetugas, &wo.KepuasanPelangganHargaLayanan,
			&wo.KepuasanPelangganTemuan, &wo.KepuasanPelangganPelanggaran,
			&wo.KepuasanPelangganHargaPelanggan, &wo.KepuasanPelangganKeterangan,
			&wo.KepuasanPelangganTdkBersedia, &wo.KepuasanPelangganAlasanGagal,
			&wo.KepuasanPelangganBaKlarifikasi, &wo.KepuasanPelangganBaKeterangan,
			&wo.KepuasanPelangganStatus, &wo.DispatchCounter, &wo.DispatchStatus, &wo.DispatchLogTime,
			&wo.CreatedDate, &wo.NoAgenda, &wo.IsMarketplace, &wo.Konfirmasi, &wo.TglKonfirmasi,
			&wo.LaporUlang, &wo.WaktuKunjungan, &wo.SendWhatsapp, &wo.IsConfirmed, &wo.CreatedBy,
			&wo.CreatedAt, &wo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		workOrders = append(workOrders, wo)
	}
	return workOrders, rows.Err()
}

func (r *WorkOrderRepositoryImpl) GetPending(ctx context.Context) ([]*entity.WorkOrder, error) {
	query := `
		SELECT * FROM work_orders 
		WHERE status < 0 OR (status IS NULL AND tgl_selesai IS NULL)
		ORDER BY tgl_kunjungan ASC, created_date ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workOrders []*entity.WorkOrder
	for rows.Next() {
		wo := &entity.WorkOrder{}
		err := rows.Scan(
			&wo.IDOrder, &wo.IDPelanggan, &wo.IDPetugas, &wo.IDKota, &wo.IDProv,
			&wo.TglOrder, &wo.TglDiterima, &wo.TglJalan, &wo.TglDiLokasi, &wo.TglPengecekan,
			&wo.TglSelesaiPengecekan, &wo.TglPembayaran, &wo.TglDikerjakan, &wo.TglSelesai,
			&wo.NamaPelanggan, &wo.NoTelp, &wo.KwhMeter, &wo.StatusMcb, &wo.StatusPadam,
			&wo.Alamat, &wo.DetailAlamat, &wo.TglKunjungan, &wo.Pengaduan, &wo.Catatan,
			&wo.LatPelanggan, &wo.LongPelanggan, &wo.BiayaTransport, &wo.BiayaPemeriksaan,
			&wo.JumlahRating, &wo.TglRating, &wo.KomentarRating, &wo.MetodeBayar, &wo.Tipe, &wo.Status,
			&wo.StatusTransfer, &wo.PenyebabBatal, &wo.Keterangan, &wo.KdUnit, &wo.KodePos,
			&wo.SumberWo, &wo.LastCallback, &wo.KetCallback, &wo.IsBerulang, &wo.FlagPayment,
			&wo.SentInvoice, &wo.WoApkt, &wo.SourceWo, &wo.Email, &wo.NamaProvinsi,
			&wo.KetCallbackDtl, &wo.AppVersion, &wo.Autodispatch, &wo.KetDispatch, &wo.ClearTemper,
			&wo.Radius, &wo.Response, &wo.Recovery, &wo.SentInvoiceBck1707, &wo.BaCancel,
			&wo.FlagBagiHasil, &wo.FlagPotongDeposit, &wo.FlagTambahDompet, &wo.IsPaidMitra,
			&wo.CallbackBy, &wo.KepuasanPelangganResponLayanan, &wo.KepuasanPelangganKualitasPetugas,
			&wo.KepuasanPelangganEtikaPetugas, &wo.KepuasanPelangganHargaLayanan,
			&wo.KepuasanPelangganTemuan, &wo.KepuasanPelangganPelanggaran,
			&wo.KepuasanPelangganHargaPelanggan, &wo.KepuasanPelangganKeterangan,
			&wo.KepuasanPelangganTdkBersedia, &wo.KepuasanPelangganAlasanGagal,
			&wo.KepuasanPelangganBaKlarifikasi, &wo.KepuasanPelangganBaKeterangan,
			&wo.KepuasanPelangganStatus, &wo.DispatchCounter, &wo.DispatchStatus, &wo.DispatchLogTime,
			&wo.CreatedDate, &wo.NoAgenda, &wo.IsMarketplace, &wo.Konfirmasi, &wo.TglKonfirmasi,
			&wo.LaporUlang, &wo.WaktuKunjungan, &wo.SendWhatsapp, &wo.IsConfirmed, &wo.CreatedBy,
			&wo.CreatedAt, &wo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		workOrders = append(workOrders, wo)
	}
	return workOrders, rows.Err()
}

func (r *WorkOrderRepositoryImpl) GetPendingByNoTelp(ctx context.Context, noTelp string) ([]*entity.WorkOrder, error) {
	query := `
		SELECT * FROM work_orders 
		WHERE no_telp = $1 AND (status < 0 OR (status IS NULL AND tgl_selesai IS NULL))
		ORDER BY tgl_kunjungan ASC
	`
	rows, err := r.db.QueryContext(ctx, query, noTelp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workOrders []*entity.WorkOrder
	for rows.Next() {
		wo := &entity.WorkOrder{}
		err := rows.Scan(
			&wo.IDOrder, &wo.IDPelanggan, &wo.IDPetugas, &wo.IDKota, &wo.IDProv,
			&wo.TglOrder, &wo.TglDiterima, &wo.TglJalan, &wo.TglDiLokasi, &wo.TglPengecekan,
			&wo.TglSelesaiPengecekan, &wo.TglPembayaran, &wo.TglDikerjakan, &wo.TglSelesai,
			&wo.NamaPelanggan, &wo.NoTelp, &wo.KwhMeter, &wo.StatusMcb, &wo.StatusPadam,
			&wo.Alamat, &wo.DetailAlamat, &wo.TglKunjungan, &wo.Pengaduan, &wo.Catatan,
			&wo.LatPelanggan, &wo.LongPelanggan, &wo.BiayaTransport, &wo.BiayaPemeriksaan,
			&wo.JumlahRating, &wo.TglRating, &wo.KomentarRating, &wo.MetodeBayar, &wo.Tipe, &wo.Status,
			&wo.StatusTransfer, &wo.PenyebabBatal, &wo.Keterangan, &wo.KdUnit, &wo.KodePos,
			&wo.SumberWo, &wo.LastCallback, &wo.KetCallback, &wo.IsBerulang, &wo.FlagPayment,
			&wo.SentInvoice, &wo.WoApkt, &wo.SourceWo, &wo.Email, &wo.NamaProvinsi,
			&wo.KetCallbackDtl, &wo.AppVersion, &wo.Autodispatch, &wo.KetDispatch, &wo.ClearTemper,
			&wo.Radius, &wo.Response, &wo.Recovery, &wo.SentInvoiceBck1707, &wo.BaCancel,
			&wo.FlagBagiHasil, &wo.FlagPotongDeposit, &wo.FlagTambahDompet, &wo.IsPaidMitra,
			&wo.CallbackBy, &wo.KepuasanPelangganResponLayanan, &wo.KepuasanPelangganKualitasPetugas,
			&wo.KepuasanPelangganEtikaPetugas, &wo.KepuasanPelangganHargaLayanan,
			&wo.KepuasanPelangganTemuan, &wo.KepuasanPelangganPelanggaran,
			&wo.KepuasanPelangganHargaPelanggan, &wo.KepuasanPelangganKeterangan,
			&wo.KepuasanPelangganTdkBersedia, &wo.KepuasanPelangganAlasanGagal,
			&wo.KepuasanPelangganBaKlarifikasi, &wo.KepuasanPelangganBaKeterangan,
			&wo.KepuasanPelangganStatus, &wo.DispatchCounter, &wo.DispatchStatus, &wo.DispatchLogTime,
			&wo.CreatedDate, &wo.NoAgenda, &wo.IsMarketplace, &wo.Konfirmasi, &wo.TglKonfirmasi,
			&wo.LaporUlang, &wo.WaktuKunjungan, &wo.SendWhatsapp, &wo.IsConfirmed, &wo.CreatedBy,
			&wo.CreatedAt, &wo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		workOrders = append(workOrders, wo)
	}
	return workOrders, rows.Err()
}

func (r *WorkOrderRepositoryImpl) Update(ctx context.Context, workOrder *entity.WorkOrder) error {
	query := `
		UPDATE work_orders SET
			id_pelanggan = $2, id_petugas = $3, id_kota = $4, id_prov = $5,
			tgl_order = $6, tgl_diterima = $7, tgl_jalan = $8, tgl_di_lokasi = $9, tgl_pengecekan = $10,
			tgl_selesai_pengecekan = $11, tgl_pembayaran = $12, tgl_dikerjakan = $13, tgl_selesai = $14,
			nama_pelanggan = $15, no_telp = $16, kwh_meter = $17, status_mcb = $18, status_padam = $19,
			alamat = $20, detail_alamat = $21, tgl_kunjungan = $22, pengaduan = $23, catatan = $24,
			lat_pelanggan = $25, long_pelanggan = $26, biaya_transport = $27, biaya_pemeriksaan = $28,
			jumlah_rating = $29, tgl_rating = $30, komentar_rating = $31, metode_bayar = $32, tipe = $33, status = $34,
			status_transfer = $35, penyebab_batal = $36, keterangan = $37, kd_unit = $38, kode_pos = $39,
			sumber_wo = $40, last_callback = $41, ket_callback = $42, is_berulang = $43, flag_payment = $44,
			sent_invoice = $45, wo_apkt = $46, source_wo = $47, email = $48, nama_provinsi = $49,
			ket_callback_dtl = $50, app_version = $51, autodispatch = $52, ket_dispatch = $53, clear_temper = $54,
			radius = $55, response = $56, recovery = $57, sent_invoice_bck1707 = $58, ba_cancel = $59,
			flag_bagi_hasil = $60, flag_potong_deposit = $61, flag_tambah_dompet = $62, is_paid_mitra = $63,
			callback_by = $64, kepuasan_pelanggan_respon_layanan = $65, kepuasan_pelanggan_kualitas_petugas = $66,
			kepuasan_pelanggan_etika_petugas = $67, kepuasan_pelanggan_harga_layanan = $68,
			kepuasan_pelanggan_temuan = $69, kepuasan_pelanggan_pelanggaran = $70,
			kepuasan_pelanggan_harga_pelanggan = $71, kepuasan_pelanggan_keterangan = $72,
			kepuasan_pelanggan_tdk_bersedia = $73, kepuasan_pelanggan_alasan_gagal = $74,
			kepuasan_pelanggan_ba_klarifikasi = $75, kepuasan_pelanggan_ba_keterangan = $76,
			kepuasan_pelanggan_status = $77, dispatch_counter = $78, dispatch_status = $79, dispatch_log_time = $80,
			created_date = $81, no_agenda = $82, is_marketplace = $83, konfirmasi = $84, tgl_konfirmasi = $85,
			lapor_ulang = $86, waktu_kunjungan = $87, send_whatsapp = $88, is_confirmed = $89, created_by = $90,
			updated_at = CURRENT_TIMESTAMP
		WHERE id_order = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		workOrder.IDOrder, workOrder.IDPelanggan, workOrder.IDPetugas, workOrder.IDKota, workOrder.IDProv,
		workOrder.TglOrder, workOrder.TglDiterima, workOrder.TglJalan, workOrder.TglDiLokasi, workOrder.TglPengecekan,
		workOrder.TglSelesaiPengecekan, workOrder.TglPembayaran, workOrder.TglDikerjakan, workOrder.TglSelesai,
		workOrder.NamaPelanggan, workOrder.NoTelp, workOrder.KwhMeter, workOrder.StatusMcb, workOrder.StatusPadam,
		workOrder.Alamat, workOrder.DetailAlamat, workOrder.TglKunjungan, workOrder.Pengaduan, workOrder.Catatan,
		workOrder.LatPelanggan, workOrder.LongPelanggan, workOrder.BiayaTransport, workOrder.BiayaPemeriksaan,
		workOrder.JumlahRating, workOrder.TglRating, workOrder.KomentarRating, workOrder.MetodeBayar, workOrder.Tipe, workOrder.Status,
		workOrder.StatusTransfer, workOrder.PenyebabBatal, workOrder.Keterangan, workOrder.KdUnit, workOrder.KodePos,
		workOrder.SumberWo, workOrder.LastCallback, workOrder.KetCallback, workOrder.IsBerulang, workOrder.FlagPayment,
		workOrder.SentInvoice, workOrder.WoApkt, workOrder.SourceWo, workOrder.Email, workOrder.NamaProvinsi,
		workOrder.KetCallbackDtl, workOrder.AppVersion, workOrder.Autodispatch, workOrder.KetDispatch, workOrder.ClearTemper,
		workOrder.Radius, workOrder.Response, workOrder.Recovery, workOrder.SentInvoiceBck1707, workOrder.BaCancel,
		workOrder.FlagBagiHasil, workOrder.FlagPotongDeposit, workOrder.FlagTambahDompet, workOrder.IsPaidMitra,
		workOrder.CallbackBy, workOrder.KepuasanPelangganResponLayanan, workOrder.KepuasanPelangganKualitasPetugas,
		workOrder.KepuasanPelangganEtikaPetugas, workOrder.KepuasanPelangganHargaLayanan,
		workOrder.KepuasanPelangganTemuan, workOrder.KepuasanPelangganPelanggaran,
		workOrder.KepuasanPelangganHargaPelanggan, workOrder.KepuasanPelangganKeterangan,
		workOrder.KepuasanPelangganTdkBersedia, workOrder.KepuasanPelangganAlasanGagal,
		workOrder.KepuasanPelangganBaKlarifikasi, workOrder.KepuasanPelangganBaKeterangan,
		workOrder.KepuasanPelangganStatus, workOrder.DispatchCounter, workOrder.DispatchStatus, workOrder.DispatchLogTime,
		workOrder.CreatedDate, workOrder.NoAgenda, workOrder.IsMarketplace, workOrder.Konfirmasi, workOrder.TglKonfirmasi,
		workOrder.LaporUlang, workOrder.WaktuKunjungan, workOrder.SendWhatsapp, workOrder.IsConfirmed, workOrder.CreatedBy,
	)
	return err
}

func (r *WorkOrderRepositoryImpl) UpdateSendWhatsapp(ctx context.Context, idOrder string, sent bool) error {
	query := `UPDATE work_orders SET send_whatsapp = $1, updated_at = CURRENT_TIMESTAMP WHERE id_order = $2`
	_, err := r.db.ExecContext(ctx, query, sent, idOrder)
	return err
}

func (r *WorkOrderRepositoryImpl) Delete(ctx context.Context, idOrder string) error {
	query := `DELETE FROM work_orders WHERE id_order = $1`
	_, err := r.db.ExecContext(ctx, query, idOrder)
	return err
}

func (r *WorkOrderRepositoryImpl) BulkCreate(ctx context.Context, workOrders []*entity.WorkOrder) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, wo := range workOrders {
		if err := r.Create(ctx, wo); err != nil {
			return fmt.Errorf("failed to create work order %s: %w", wo.IDOrder, err)
		}
	}

	return tx.Commit()
}
