package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"smart_alert_system/internal/config"
	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/infrastructure/database"
	infraRepo "smart_alert_system/internal/infrastructure/repository"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <csv_file_path>", os.Args[0])
	}

	csvFile := os.Args[1]

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Open CSV file
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)
	reader.Comma = '|' // Use pipe as delimiter
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read header: %v", err)
	}
	log.Printf("Header: %v", header)

	// Create repository
	workOrderRepo := infraRepo.NewWorkOrderRepository(db.DB)

	ctx := context.Background()
	successCount := 0
	errorCount := 0

	// Read and process each row
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading record: %v", err)
			errorCount++
			continue
		}

		// Parse work order from CSV record
		workOrder, err := parseWorkOrderFromCSV(header, record)
		if err != nil {
			log.Printf("Error parsing work order: %v", err)
			errorCount++
			continue
		}

		// Insert into database
		if err := workOrderRepo.Create(ctx, workOrder); err != nil {
			log.Printf("Error inserting work order %s: %v", workOrder.IDOrder, err)
			errorCount++
			continue
		}

		successCount++
		if successCount%100 == 0 {
			log.Printf("Processed %d work orders...", successCount)
		}
	}

	log.Printf("Import completed: %d successful, %d errors", successCount, errorCount)
}

func parseWorkOrderFromCSV(header, record []string) (*entity.WorkOrder, error) {
	if len(record) != len(header) {
		return nil, fmt.Errorf("record length (%d) doesn't match header length (%d)", len(record), len(header))
	}

	wo := &entity.WorkOrder{}

	// Create map for easy lookup
	recordMap := make(map[string]string)
	for i, h := range header {
		if i < len(record) {
			recordMap[strings.TrimSpace(h)] = strings.TrimSpace(record[i])
		}
	}

	// Parse fields
	wo.IDOrder = getString(recordMap, "id_order")
	wo.IDPelanggan = getNullString(recordMap, "id_pelanggan")
	wo.IDPetugas = getNullString(recordMap, "id_petugas")
	wo.IDKota = getNullString(recordMap, "id_kota")
	wo.IDProv = getNullString(recordMap, "id_prov")
	wo.TglOrder = getNullTime(recordMap, "tgl_order")
	wo.TglDiterima = getNullTime(recordMap, "tgl_diterima")
	wo.TglJalan = getNullTime(recordMap, "tgl_jalan")
	wo.TglDiLokasi = getNullTime(recordMap, "tgl_di_lokasi")
	wo.TglPengecekan = getNullTime(recordMap, "tgl_pengecekan")
	wo.TglSelesaiPengecekan = getNullTime(recordMap, "tgl_selesai_pengecekan")
	wo.TglPembayaran = getNullTime(recordMap, "tgl_pembayaran")
	wo.TglDikerjakan = getNullTime(recordMap, "tgl_dikerjakan")
	wo.TglSelesai = getNullTime(recordMap, "tgl_selesai")
	wo.NamaPelanggan = getNullString(recordMap, "nama_pelanggan")
	wo.NoTelp = getNullString(recordMap, "no_telp")
	wo.KwhMeter = getNullString(recordMap, "kwh_meter")
	wo.StatusMcb = getNullString(recordMap, "status_mcb")
	wo.StatusPadam = getNullString(recordMap, "status_padam")
	wo.Alamat = getNullString(recordMap, "alamat")
	wo.DetailAlamat = getNullString(recordMap, "detail_alamat")
	wo.TglKunjungan = getNullTime(recordMap, "tgl_kunjungan")
	wo.Pengaduan = getNullString(recordMap, "pengaduan")
	wo.Catatan = getNullString(recordMap, "catatan")
	wo.LatPelanggan = getNullFloat64(recordMap, "lat_pelanggan")
	wo.LongPelanggan = getNullFloat64(recordMap, "long_pelanggan")
	wo.BiayaTransport = getNullFloat64(recordMap, "biaya_transport")
	wo.BiayaPemeriksaan = getNullFloat64(recordMap, "biaya_pemeriksaan")
	wo.JumlahRating = getNullInt64(recordMap, "jumlah_rating")
	wo.TglRating = getNullTime(recordMap, "tgl_rating")
	wo.KomentarRating = getNullString(recordMap, "komentar_rating")
	wo.MetodeBayar = getNullString(recordMap, "metode_bayar")
	wo.Tipe = getNullString(recordMap, "tipe")
	wo.Status = getNullInt64(recordMap, "status")
	wo.StatusTransfer = getBool(recordMap, "status_transfer")
	wo.PenyebabBatal = getNullString(recordMap, "penyebab_batal")
	wo.Keterangan = getNullString(recordMap, "keterangan")
	wo.KdUnit = getNullString(recordMap, "kd_unit")
	wo.KodePos = getNullString(recordMap, "kode_pos")
	wo.SumberWo = getNullString(recordMap, "sumber_wo")
	wo.LastCallback = getNullTime(recordMap, "last_callback")
	wo.KetCallback = getNullString(recordMap, "ket_callback")
	wo.IsBerulang = getBool(recordMap, "is_berulang")
	wo.FlagPayment = getBool(recordMap, "flag_payment")
	wo.SentInvoice = getBool(recordMap, "sent_invoice")
	wo.WoApkt = getNullString(recordMap, "wo_apkt")
	wo.SourceWo = getNullString(recordMap, "source_wo")
	wo.Email = getNullString(recordMap, "email")
	wo.NamaProvinsi = getNullString(recordMap, "nama_provinsi")
	wo.KetCallbackDtl = getNullString(recordMap, "ket_callback_dtl")
	wo.AppVersion = getNullString(recordMap, "app_version")
	wo.Autodispatch = getInt(recordMap, "autodispatch")
	wo.KetDispatch = getNullString(recordMap, "ket_dispatch")
	wo.ClearTemper = getNullString(recordMap, "clear_temper")
	wo.Radius = getNullFloat64(recordMap, "radius")
	wo.Response = getNullString(recordMap, "response")
	wo.Recovery = getNullString(recordMap, "recovery")
	wo.SentInvoiceBck1707 = getBool(recordMap, "sent_invoice_bck1707")
	wo.BaCancel = getNullString(recordMap, "ba_cancel")
	wo.FlagBagiHasil = getBool(recordMap, "flag_bagi_hasil")
	wo.FlagPotongDeposit = getBool(recordMap, "flag_potong_deposit")
	wo.FlagTambahDompet = getBool(recordMap, "flag_tambah_dompet")
	wo.IsPaidMitra = getBool(recordMap, "is_paid_mitra")
	wo.CallbackBy = getNullString(recordMap, "callback_by")
	wo.KepuasanPelangganResponLayanan = getNullString(recordMap, "kepuasan_pelanggan_respon_layanan")
	wo.KepuasanPelangganKualitasPetugas = getNullString(recordMap, "kepuasan_pelanggan_kualitas_petugas")
	wo.KepuasanPelangganEtikaPetugas = getNullString(recordMap, "kepuasan_pelanggan_etika_petugas")
	wo.KepuasanPelangganHargaLayanan = getNullString(recordMap, "kepuasan_pelanggan_harga_layanan")
	wo.KepuasanPelangganTemuan = getNullString(recordMap, "kepuasan_pelanggan_temuan")
	wo.KepuasanPelangganPelanggaran = getNullString(recordMap, "kepuasan_pelanggan_pelanggaran")
	wo.KepuasanPelangganHargaPelanggan = getNullString(recordMap, "kepuasan_pelanggan_harga_pelanggan")
	wo.KepuasanPelangganKeterangan = getNullString(recordMap, "kepuasan_pelanggan_keterangan")
	wo.KepuasanPelangganTdkBersedia = getNullString(recordMap, "kepuasan_pelanggan_tdk_bersedia")
	wo.KepuasanPelangganAlasanGagal = getNullString(recordMap, "kepuasan_pelanggan_alasan_gagal")
	wo.KepuasanPelangganBaKlarifikasi = getNullString(recordMap, "kepuasan_pelanggan_ba_klarifikasi")
	wo.KepuasanPelangganBaKeterangan = getNullString(recordMap, "kepuasan_pelanggan_ba_keterangan")
	wo.KepuasanPelangganStatus = getNullString(recordMap, "kepuasan_pelanggan_status")
	wo.DispatchCounter = getInt(recordMap, "dispatch_counter")
	wo.DispatchStatus = getNullString(recordMap, "dispatch_status")
	wo.DispatchLogTime = getNullTime(recordMap, "dispatch_log_time")
	wo.CreatedDate = getTime(recordMap, "created_date", time.Now())
	wo.NoAgenda = getNullString(recordMap, "no_agenda")
	wo.IsMarketplace = getBool(recordMap, "is_marketplace")
	wo.Konfirmasi = getNullString(recordMap, "konfirmasi")
	wo.TglKonfirmasi = getNullTime(recordMap, "tgl_konfirmasi")
	wo.LaporUlang = getInt(recordMap, "lapor_ulang")
	wo.WaktuKunjungan = getNullTime(recordMap, "waktu_kunjungan")
	wo.SendWhatsapp = getBool(recordMap, "send_whatsapp")
	wo.IsConfirmed = getBool(recordMap, "is_confirmed")
	wo.CreatedBy = getNullString(recordMap, "created_by")

	return wo, nil
}

// Helper functions for parsing
func getString(m map[string]string, key string) string {
	if val, ok := m[key]; ok && val != "" {
		return val
	}
	return ""
}

func getNullString(m map[string]string, key string) sql.NullString {
	if val, ok := m[key]; ok && val != "" {
		return sql.NullString{String: val, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func getNullTime(m map[string]string, key string) sql.NullTime {
	if val, ok := m[key]; ok && val != "" {
		// Try multiple time formats
		formats := []string{
			"2006-01-02 15:04:05.000",
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05.000Z",
			"2006-01-02",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, val); err == nil {
				return sql.NullTime{Time: t, Valid: true}
			}
		}
	}
	return sql.NullTime{Valid: false}
}

func getTime(m map[string]string, key string, defaultValue time.Time) time.Time {
	if val, ok := m[key]; ok && val != "" {
		formats := []string{
			"2006-01-02 15:04:05.000",
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05.000Z",
			"2006-01-02",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, val); err == nil {
				return t
			}
		}
	}
	return defaultValue
}

func getNullFloat64(m map[string]string, key string) sql.NullFloat64 {
	if val, ok := m[key]; ok && val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return sql.NullFloat64{Float64: f, Valid: true}
		}
	}
	return sql.NullFloat64{Valid: false}
}

func getNullInt64(m map[string]string, key string) sql.NullInt64 {
	if val, ok := m[key]; ok && val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return sql.NullInt64{Int64: i, Valid: true}
		}
	}
	return sql.NullInt64{Valid: false}
}

func getInt(m map[string]string, key string) int {
	if val, ok := m[key]; ok && val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return 0
}

func getBool(m map[string]string, key string) bool {
	if val, ok := m[key]; ok && val != "" {
		val = strings.ToLower(strings.TrimSpace(val))
		return val == "true" || val == "1" || val == "yes" || val == "t"
	}
	return false
}
