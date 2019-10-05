package hansard

import (
	"testing"
)

func Test_isStartOfQuestionSection(t *testing.T) {
	type args struct {
		rowContents []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"happy #1", args{[]string{"PERTANYAAN-PERTANYAAN MULUT DARIPADA ", "Y.B. PUAN JAMALIAH BINTI JAMALUDDIN  ", "(N36 BANDAR UTAMA) ", "TAJUK : ISU KESIHATAN DAN KESELAMATAN DI KAWASAN TANAH  ", "261.  Bertanya kepada Y.A.B. Dato'' Menteri Besar:- ", "a)     ", "Sekiranya terdapat kawasan-kawasan tanah terbiar yang menyumbang kepada "}}, true},
		{"happy #2", args{[]string{"PERTANYAAN-PERTANYAAN BERTULIS DARIPADA ", "Y.B. PUAN LEE KEE HIONG  ", "(N06 KUALA KUBU BAHARU) ", "TAJUK : TANAH PERBADANAN KEMAJUAN PERTANIAN SELANGOR (PKPS) ", "57.   Bertanya kepada Y.A.B. Dato'' Menteri Besar:- ", "a)         Senaraikan tanah milikan PKPS yang diusahakan secara sendiri atau usaha ", " supaya 30.00 ratio tanah 30. Saya bertanya adakah in benar?", "7 Ulu Tinggi 9.4157 Ternakan ayam PKPS "}}, true},
		{"happy #3", args{[]string{"PERTANYAAN-PERTANYAAN MULUT DARIPADA ", "Y.B. PUAN DR. DAROYAH BINTI ALWI  ", "(N43 SEMENTA) ", "TAJUK : PERUNTUKAN UNTUK MEMBERDAYAKAN PUSAT WANITA ", "1.  Bertanya kepada Y.A.B. Dato'' Menteri Besar:- ", "Sila nyatakan berapa banyak program PWB telah dilaksanakan sepanjang ", "tahun 2018, mengikut pecahan kategori program. ", "Berapakah peratus peruntukan PWB telah di belanjakan untuk tahun 2018, dan "}}, true},
		{"happy #4", args{[]string{"� Perayaan Pesta Ponggal Peringkat Negeri Selangor  ", "PERTANYAAN-PERTANYAAN BERTULIS DARIPADA ", "Y.B. PUAN JUWAIRIYA BINTI ZULKIFLI  ", "(N10 BUKIT MELAWATI) ", "TAJUK : MUZIUM PERMAINAN TRADISIONAL RAKYAT DI BUKIT MALAWATI ", "8.   Bertanya kepada Y.A.B. Dato'' Menteri Besar:- ", "a)         Apakah status muzium permainan tradisional\n rakyat di Bukit Malawati \nkerana ", "menyediakan gunatenaga dan pendapatan Negeri Selangor. Berdasarkan ", "ii. Pembangunan semula Kesultanan Awal Negeri Selangor dan ", "Selangor yang dibuka seawal tahun 1700an. Ianya berpusat di Bukit ", "Jugra yang telah dibuka kepada pelawat mulai tahun 2004 yang lalu. ", "Almathum Sultan Ala’edin Sulaiman Shah. Walau� bagaimanapun 3,000 orang"}}, true},
		{"sad #1", args{[]string{"JAWAPAN: ", "a) Untuk bantuan membaikpulih bangunan, bumbung, cat semula dinding dan ", "keperluan lain untuk keceriaan Pangsapuri Sri Sementa terdapat dua pilihan ", "s ", "2017 125,757.12 4.12 30.07 ", "b) Data pengeluaran bahan batuan dan mineral bagi tahun 2017 dan 2018 adalah "}}, false},
		{"sad #2", args{[]string{"BIL DAERAH MUKIM ", "LUAS ", "(HEKTAR) ", "7 Ulu Tinggi 9.4157 Ternakan ayam PKPS ", "8 ", "Ulu Yam 80.39 Ladang Kelapa Sawit PKPS   ", "9 Kerling 374.9 ", "Kerling PKPS "}}, false},
		//{"", args{[]string{"", ""}}, true},
		//{"", args{[]string{"", ""}}, true},
		//{"", args{[]string{"", ""}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pdfPage := PDFPage{PDFTxtSameLines: tt.args.rowContents}
			if got := isStartOfQuestionSection(pdfPage); got != tt.want {
				t.Errorf("isStartOfQuestionSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractQuestionNum(t *testing.T) {
	type args struct {
		rowContent string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"happy #1", args{"261.  Bertanya kepada Y.A.B. Dato'' Menteri Besar:- "}, "261", false},
		{"happy #2", args{" � 50 bertanya kepada yab menteri besar Azmin ALI "}, "50", false},
		{"happy odd #3", args{"  01)   Bertanya kepada Y.A.B. "}, "01", false},
		{"happy #4", args{"41.   Bertanya kepada Y.A.B. Dato'' Menteri Besar:- "}, "41", false},
		{"happy #5", args{"  3 # BERTanya KEPADA Menteri Pertanian"}, "3", false},
		{"sad #1", args{"TAJUK : PERUNTUKAN CERIA 2019 "}, "", false},
		{"sad #2", args{"(N43 SEMENTA) "}, "", false},
		//{"break #1", args{"\n\n\n 50 bertanya kepada yab menteri besar Azmin ALI "}, "50", false},
		//{"", args{""}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractQuestionNum(tt.args.rowContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractQuestionNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractQuestionNum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_detectHansardType(t *testing.T) {
	type args struct {
		firstPage PDFPage
	}
	tests := []struct {
		name string
		args args
		want HansardType
	}{
		{"happy #1", args{PDFPage{PDFTxtSameLines: []string{"PERTANYAAN-PERTANYAAN MULUT DARIPADA ", "Y.B. PUAN LEE KEE HIONG ", "(N06 KUALA KUBU BAHARU) ", "TAJUK : PROGRAM HARI TANPA KENDERAAN ", "4.  Bertanya kepada Y.A.B. Dato'' Menteri Besar:- ", "a)   "}}}, HANSARD_SPOKEN},
		{"happy #2", args{PDFPage{PDFTxtSameLines: []string{"PERTANYAAN-PERTANYAAN BERTULIS DARIPADA ", "(N43 SEMENTA) ", "TAJUK : BANTUAN SKIM CERIA", "'2.   Bertanya kepada Y.A.B. Dato'' Menteri Besar:- "}}}, HANSARD_WRITTEN},
		{"not quite #1", args{PDFPage{PDFTxtSameLines: []string{"YAB Yeo", "PerTANYAAN mengenai Sampah yang dibuang"}}}, -1},
		{"sad #1", args{PDFPage{PDFTxtSameLines: []string{"Jumlah Kutipan Royalti (RM) ", "2017 80,127,921.02 Juta 125,757.12 ", "Di Bawah Tanah Kerajaan dan Tanah Hakmilik  "}}}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := detectHansardType(tt.args.firstPage); got != tt.want {
				t.Errorf("detectHansardType() = %v, want %v", got, tt.want)
			}
		})
	}
}
