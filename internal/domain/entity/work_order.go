package entity

import (
	"database/sql"
	"time"
)

type WorkOrder struct {
	IDOrder                        string         `json:"id_order" db:"id_order"`
	IDPelanggan                    sql.NullString `json:"id_pelanggan" db:"id_pelanggan"`
	IDPetugas                      sql.NullString `json:"id_petugas" db:"id_petugas"`
	IDKota                         sql.NullString `json:"id_kota" db:"id_kota"`
	IDProv                         sql.NullString `json:"id_prov" db:"id_prov"`
	TglOrder                       sql.NullTime   `json:"tgl_order" db:"tgl_order"`
	TglDiterima                    sql.NullTime   `json:"tgl_diterima" db:"tgl_diterima"`
	TglJalan                       sql.NullTime   `json:"tgl_jalan" db:"tgl_jalan"`
	TglDiLokasi                     sql.NullTime   `json:"tgl_di_lokasi" db:"tgl_di_lokasi"`
	TglPengecekan                   sql.NullTime   `json:"tgl_pengecekan" db:"tgl_pengecekan"`
	TglSelesaiPengecekan           sql.NullTime   `json:"tgl_selesai_pengecekan" db:"tgl_selesai_pengecekan"`
	TglPembayaran                  sql.NullTime   `json:"tgl_pembayaran" db:"tgl_pembayaran"`
	TglDikerjakan                  sql.NullTime   `json:"tgl_dikerjakan" db:"tgl_dikerjakan"`
	TglSelesai                     sql.NullTime   `json:"tgl_selesai" db:"tgl_selesai"`
	NamaPelanggan                  sql.NullString `json:"nama_pelanggan" db:"nama_pelanggan"`
	NoTelp                         sql.NullString `json:"no_telp" db:"no_telp"`
	KwhMeter                       sql.NullString `json:"kwh_meter" db:"kwh_meter"`
	StatusMcb                      sql.NullString `json:"status_mcb" db:"status_mcb"`
	StatusPadam                    sql.NullString `json:"status_padam" db:"status_padam"`
	Alamat                         sql.NullString `json:"alamat" db:"alamat"`
	DetailAlamat                   sql.NullString `json:"detail_alamat" db:"detail_alamat"`
	TglKunjungan                   sql.NullTime   `json:"tgl_kunjungan" db:"tgl_kunjungan"`
	Pengaduan                      sql.NullString `json:"pengaduan" db:"pengaduan"`
	Catatan                        sql.NullString `json:"catatan" db:"catatan"`
	LatPelanggan                   sql.NullFloat64 `json:"lat_pelanggan" db:"lat_pelanggan"`
	LongPelanggan                  sql.NullFloat64 `json:"long_pelanggan" db:"long_pelanggan"`
	BiayaTransport                 sql.NullFloat64 `json:"biaya_transport" db:"biaya_transport"`
	BiayaPemeriksaan               sql.NullFloat64 `json:"biaya_pemeriksaan" db:"biaya_pemeriksaan"`
	JumlahRating                   sql.NullInt64   `json:"jumlah_rating" db:"jumlah_rating"`
	TglRating                      sql.NullTime   `json:"tgl_rating" db:"tgl_rating"`
	KomentarRating                 sql.NullString `json:"komentar_rating" db:"komentar_rating"`
	MetodeBayar                    sql.NullString `json:"metode_bayar" db:"metode_bayar"`
	Tipe                           sql.NullString `json:"tipe" db:"tipe"`
	Status                         sql.NullInt64  `json:"status" db:"status"`
	StatusTransfer                 bool            `json:"status_transfer" db:"status_transfer"`
	PenyebabBatal                  sql.NullString `json:"penyebab_batal" db:"penyebab_batal"`
	Keterangan                     sql.NullString `json:"keterangan" db:"keterangan"`
	KdUnit                         sql.NullString `json:"kd_unit" db:"kd_unit"`
	KodePos                        sql.NullString `json:"kode_pos" db:"kode_pos"`
	SumberWo                       sql.NullString `json:"sumber_wo" db:"sumber_wo"`
	LastCallback                   sql.NullTime   `json:"last_callback" db:"last_callback"`
	KetCallback                    sql.NullString `json:"ket_callback" db:"ket_callback"`
	IsBerulang                     bool            `json:"is_berulang" db:"is_berulang"`
	FlagPayment                    bool            `json:"flag_payment" db:"flag_payment"`
	SentInvoice                    bool            `json:"sent_invoice" db:"sent_invoice"`
	WoApkt                         sql.NullString `json:"wo_apkt" db:"wo_apkt"`
	SourceWo                       sql.NullString `json:"source_wo" db:"source_wo"`
	Email                          sql.NullString `json:"email" db:"email"`
	NamaProvinsi                   sql.NullString `json:"nama_provinsi" db:"nama_provinsi"`
	KetCallbackDtl                 sql.NullString `json:"ket_callback_dtl" db:"ket_callback_dtl"`
	AppVersion                     sql.NullString `json:"app_version" db:"app_version"`
	Autodispatch                   int            `json:"autodispatch" db:"autodispatch"`
	KetDispatch                    sql.NullString `json:"ket_dispatch" db:"ket_dispatch"`
	ClearTemper                    sql.NullString `json:"clear_temper" db:"clear_temper"`
	Radius                         sql.NullFloat64 `json:"radius" db:"radius"`
	Response                       sql.NullString `json:"response" db:"response"`
	Recovery                       sql.NullString `json:"recovery" db:"recovery"`
	SentInvoiceBck1707             bool            `json:"sent_invoice_bck1707" db:"sent_invoice_bck1707"`
	BaCancel                       sql.NullString `json:"ba_cancel" db:"ba_cancel"`
	FlagBagiHasil                  bool            `json:"flag_bagi_hasil" db:"flag_bagi_hasil"`
	FlagPotongDeposit              bool            `json:"flag_potong_deposit" db:"flag_potong_deposit"`
	FlagTambahDompet               bool            `json:"flag_tambah_dompet" db:"flag_tambah_dompet"`
	IsPaidMitra                    bool            `json:"is_paid_mitra" db:"is_paid_mitra"`
	CallbackBy                     sql.NullString `json:"callback_by" db:"callback_by"`
	KepuasanPelangganResponLayanan sql.NullString `json:"kepuasan_pelanggan_respon_layanan" db:"kepuasan_pelanggan_respon_layanan"`
	KepuasanPelangganKualitasPetugas sql.NullString `json:"kepuasan_pelanggan_kualitas_petugas" db:"kepuasan_pelanggan_kualitas_petugas"`
	KepuasanPelangganEtikaPetugas  sql.NullString `json:"kepuasan_pelanggan_etika_petugas" db:"kepuasan_pelanggan_etika_petugas"`
	KepuasanPelangganHargaLayanan sql.NullString `json:"kepuasan_pelanggan_harga_layanan" db:"kepuasan_pelanggan_harga_layanan"`
	KepuasanPelangganTemuan        sql.NullString `json:"kepuasan_pelanggan_temuan" db:"kepuasan_pelanggan_temuan"`
	KepuasanPelangganPelanggaran  sql.NullString `json:"kepuasan_pelanggan_pelanggaran" db:"kepuasan_pelanggan_pelanggaran"`
	KepuasanPelangganHargaPelanggan sql.NullString `json:"kepuasan_pelanggan_harga_pelanggan" db:"kepuasan_pelanggan_harga_pelanggan"`
	KepuasanPelangganKeterangan    sql.NullString `json:"kepuasan_pelanggan_keterangan" db:"kepuasan_pelanggan_keterangan"`
	KepuasanPelangganTdkBersedia   sql.NullString `json:"kepuasan_pelanggan_tdk_bersedia" db:"kepuasan_pelanggan_tdk_bersedia"`
	KepuasanPelangganAlasanGagal   sql.NullString `json:"kepuasan_pelanggan_alasan_gagal" db:"kepuasan_pelanggan_alasan_gagal"`
	KepuasanPelangganBaKlarifikasi sql.NullString `json:"kepuasan_pelanggan_ba_klarifikasi" db:"kepuasan_pelanggan_ba_klarifikasi"`
	KepuasanPelangganBaKeterangan  sql.NullString `json:"kepuasan_pelanggan_ba_keterangan" db:"kepuasan_pelanggan_ba_keterangan"`
	KepuasanPelangganStatus        sql.NullString `json:"kepuasan_pelanggan_status" db:"kepuasan_pelanggan_status"`
	DispatchCounter                int            `json:"dispatch_counter" db:"dispatch_counter"`
	DispatchStatus                 sql.NullString `json:"dispatch_status" db:"dispatch_status"`
	DispatchLogTime                sql.NullTime   `json:"dispatch_log_time" db:"dispatch_log_time"`
	CreatedDate                    time.Time      `json:"created_date" db:"created_date"`
	NoAgenda                       sql.NullString `json:"no_agenda" db:"no_agenda"`
	IsMarketplace                  bool            `json:"is_marketplace" db:"is_marketplace"`
	Konfirmasi                     sql.NullString `json:"konfirmasi" db:"konfirmasi"`
	TglKonfirmasi                  sql.NullTime   `json:"tgl_konfirmasi" db:"tgl_konfirmasi"`
	LaporUlang                     int            `json:"lapor_ulang" db:"lapor_ulang"`
	WaktuKunjungan                 sql.NullTime   `json:"waktu_kunjungan" db:"waktu_kunjungan"`
	SendWhatsapp                   bool            `json:"send_whatsapp" db:"send_whatsapp"`
	IsConfirmed                    bool            `json:"is_confirmed" db:"is_confirmed"`
	CreatedBy                      sql.NullString `json:"created_by" db:"created_by"`
	CreatedAt                      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt                      time.Time      `json:"updated_at" db:"updated_at"`
}

// GetTglKunjungan returns the visit date/time, prioritizing tgl_kunjungan, then waktu_kunjungan
func (wo *WorkOrder) GetTglKunjungan() *time.Time {
	if wo.TglKunjungan.Valid {
		return &wo.TglKunjungan.Time
	}
	if wo.WaktuKunjungan.Valid {
		return &wo.WaktuKunjungan.Time
	}
	return nil
}

// GetNoTelpClean returns clean phone number (without country code prefix)
func (wo *WorkOrder) GetNoTelpClean() string {
	if !wo.NoTelp.Valid {
		return ""
	}
	noTelp := wo.NoTelp.String
	// Remove country code if present (e.g., 62 -> 0)
	if len(noTelp) > 0 && noTelp[0] == '6' && len(noTelp) > 1 && noTelp[1] == '2' {
		noTelp = "0" + noTelp[2:]
	}
	return noTelp
}

// IsPending returns true if work order is pending (status < 0 or specific pending statuses)
func (wo *WorkOrder) IsPending() bool {
	if !wo.Status.Valid {
		return false
	}
	return wo.Status.Int64 < 0
}

// IsCompleted returns true if work order is completed
func (wo *WorkOrder) IsCompleted() bool {
	if !wo.Status.Valid {
		return false
	}
	return wo.Status.Int64 > 0 && wo.TglSelesai.Valid
}

