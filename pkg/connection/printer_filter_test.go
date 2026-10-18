package connection

import (
	"runtime"
	"testing"
)

func TestInterpretStatus(t *testing.T) {
	tests := []struct {
		name       string
		status     uint32
		attributes uint32
		want       PrinterState
	}{
		{
			name:       "zero status and attributes is ready",
			status:     0,
			attributes: 0,
			want:       StateReady,
		},
		{
			name:       "status offline",
			status:     PrinterStatusOffline,
			attributes: 0,
			want:       StateOffline,
		},
		{
			name:       "status PRINTER_STATUS_OFFLINE alias",
			status:     PRINTER_STATUS_OFFLINE,
			attributes: 0,
			want:       StateOffline,
		},
		{
			name:       "attribute work offline with zero status",
			status:     0,
			attributes: PrinterAttributeWorkOffline,
			want:       StateOffline,
		},
		{
			name:       "attribute PRINTER_ATTRIBUTE_WORK_OFFLINE alias",
			status:     0,
			attributes: PRINTER_ATTRIBUTE_WORK_OFFLINE,
			want:       StateOffline,
		},
		{
			name:       "status paused",
			status:     PrinterStatusPaused,
			attributes: 0,
			want:       StatePaused,
		},
		{
			name:       "status PRINTER_STATUS_PAUSED alias",
			status:     PRINTER_STATUS_PAUSED,
			attributes: 0,
			want:       StatePaused,
		},
		{
			name:       "status error",
			status:     PrinterStatusError,
			attributes: 0,
			want:       StateError,
		},
		{
			name:       "status paper jam",
			status:     PrinterStatusPaperJam,
			attributes: 0,
			want:       StateError,
		},
		{
			name:       "status paper out",
			status:     PrinterStatusPaperOut,
			attributes: 0,
			want:       StateError,
		},
		{
			name:       "status door open",
			status:     PrinterStatusDoorOpen,
			attributes: 0,
			want:       StateError,
		},
		{
			name:       "status paper problem",
			status:     PrinterStatusPaperProblem,
			attributes: 0,
			want:       StateError,
		},
		{
			name:       "status no toner",
			status:     PrinterStatusNoToner,
			attributes: 0,
			want:       StateError,
		},
		{
			name:       "status out of memory",
			status:     PrinterStatusOutOfMemory,
			attributes: 0,
			want:       StateError,
		},
		{
			name:       "status user intervention",
			status:     PrinterStatusUserIntervention,
			attributes: 0,
			want:       StateError,
		},
		{
			name:       "status printing is ready",
			status:     PrinterStatusPrinting,
			attributes: 0,
			want:       StateReady,
		},
		{
			name:       "status processing is ready",
			status:     PrinterStatusProcessing,
			attributes: 0,
			want:       StateReady,
		},
		{
			name:       "status busy is ready",
			status:     PrinterStatusBusy,
			attributes: 0,
			want:       StateReady,
		},
		{
			name:       "offline precedence over paused",
			status:     PrinterStatusOffline | PrinterStatusPaused,
			attributes: 0,
			want:       StateOffline,
		},
		{
			name:       "offline attribute precedence over error",
			status:     PrinterStatusError,
			attributes: PrinterAttributeWorkOffline,
			want:       StateOffline,
		},
		{
			name:       "unknown status flags",
			status:     PrinterStatusPowerSave,
			attributes: 0,
			want:       StateUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := interpretStatus(tt.status, tt.attributes)
			if got != tt.want {
				t.Errorf("interpretStatus(%#x, %#x) = %v, want %v", tt.status, tt.attributes, got, tt.want)
			}
		})
	}
}

func TestIsVirtualPrinter(t *testing.T) {
	tests := []struct {
		name   string
		pName  string
		driver string
		port   string
		want   bool
	}{
		{
			name:   "Microsoft Print to PDF",
			pName:  "Microsoft Print to PDF",
			driver: "Microsoft Print To PDF",
			port:   "PORTPROMPT:",
			want:   true,
		},
		{
			name:   "Microsoft XPS Document Writer",
			pName:  "Microsoft XPS Document Writer",
			driver: "Microsoft XPS Document Writer v4",
			port:   "XPSPort:",
			want:   true,
		},
		{
			name:   "OneNote",
			pName:  "Send to OneNote 2016",
			driver: "Send to Microsoft OneNote 16 Driver",
			port:   "nul:",
			want:   true,
		},
		{
			name:   "Adobe PDF",
			pName:  "Adobe PDF",
			driver: "Adobe PDF Converter",
			port:   "Documents\\*.pdf",
			want:   true,
		},
		{
			name:   "CutePDF",
			pName:  "CutePDF Writer",
			driver: "CutePDF Writer Driver",
			port:   "CPW2:",
			want:   true,
		},
		{
			name:   "Fax",
			pName:  "Fax",
			driver: "Microsoft Shared Fax Driver",
			port:   "SHRFAX:",
			want:   true,
		},
		{
			name:   "File port virtual prefix",
			pName:  "Custom Exporter",
			driver: "Generic Driver",
			port:   "FILE:",
			want:   true,
		},
		{
			name:   "Nul port virtual prefix",
			pName:  "Null Sink",
			driver: "Raw Driver",
			port:   "NUL:",
			want:   true,
		},
		{
			name:   "Physical EPSON thermal printer",
			pName:  "EPSON TM-T88V",
			driver: "EPSON TM-T88V Receipt",
			port:   "USB001",
			want:   false,
		},
		{
			name:   "Physical HP LaserJet",
			pName:  "HP LaserJet Pro MFP M428fdw",
			driver: "HP LaserJet Pro PCL 6",
			port:   "WSD-9876",
			want:   false,
		},
		{
			name:   "Physical Star printer",
			pName:  "Star TSP100",
			driver: "Star TSP100 Cutter",
			port:   "TCP_192.168.1.100",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isVirtualPrinter(tt.pName, tt.driver, tt.port)
			if got != tt.want {
				t.Errorf("isVirtualPrinter(%q, %q, %q) = %v, want %v", tt.pName, tt.driver, tt.port, got, tt.want)
			}
		})
	}
}

func TestDetectPrinterType(t *testing.T) {
	tests := []struct {
		name       string
		port       string
		driver     string
		printer    string
		wantType   string
	}{
		{
			name:     "virtual PDF printer",
			port:     "PORTPROMPT:",
			driver:   "Microsoft Print To PDF",
			printer:  "Microsoft Print to PDF",
			wantType: "virtual",
		},
		{
			name:     "virtual XPS printer",
			port:     "XPSPort:",
			driver:   "Microsoft XPS Document Writer v4",
			printer:  "Microsoft XPS Document Writer",
			wantType: "virtual",
		},
		{
			name:     "virtual OneNote printer",
			port:     "nul:",
			driver:   "Send to Microsoft OneNote 16 Driver",
			printer:  "OneNote",
			wantType: "virtual",
		},
		{
			name:     "thermal EPSON USB",
			port:     "USB001",
			driver:   "EPSON TM-T88V Receipt",
			printer:  "Receipt Printer",
			wantType: "thermal",
		},
		{
			name:     "thermal Star Network",
			port:     "TCP_192.168.1.200",
			driver:   "Star TSP143IIIL",
			printer:  "Kitchen Star",
			wantType: "thermal",
		},
		{
			name:     "thermal POS-80 Serial",
			port:     "COM1",
			driver:   "POS-80 Series Driver",
			printer:  "POS-80",
			wantType: "thermal",
		},
		{
			name:     "thermal 58mm",
			port:     "USB002",
			driver:   "58mm Series Printer",
			printer:  "Mini Receipt",
			wantType: "thermal",
		},
		{
			name:     "thermal Citizen",
			port:     "USB003",
			driver:   "Citizen CT-S310II",
			printer:  "Front Desk",
			wantType: "thermal",
		},
		{
			name:     "thermal Bixolon",
			port:     "USB004",
			driver:   "Bixolon SRP-350plusIII",
			printer:  "Bar Ticket",
			wantType: "thermal",
		},
		{
			name:     "thermal ESC/POS pattern",
			port:     "USB005",
			driver:   "Generic ESC/POS Driver",
			printer:  "POS Terminal",
			wantType: "thermal",
		},
		{
			name:     "network UNC path non-thermal",
			port:     `\\PRINTSERVER\OfficeLaser`,
			driver:   "HP LaserJet 4000",
			printer:  "Office LaserJet",
			wantType: "network",
		},
		{
			name:     "network WSD non-thermal",
			port:     "WSD-12345",
			driver:   "Canon iR-ADV C5535",
			printer:  "Canon Copier",
			wantType: "network",
		},
		{
			name:     "network IP_ prefix non-thermal",
			port:     "IP_192.168.1.150",
			driver:   "Brother HL-L2350DW",
			printer:  "Office Brother",
			wantType: "network",
		},
		{
			name:     "network TCP prefix non-thermal",
			port:     "tcp_10.0.0.88",
			driver:   "Lexmark MS810",
			printer:  "Floor Lexmark",
			wantType: "network",
		},
		{
			name:     "usb non-thermal",
			port:     "USB001",
			driver:   "HP DeskJet 2700 series",
			printer:  "Color DeskJet",
			wantType: "usb",
		},
		{
			name:     "serial COM non-thermal",
			port:     "COM3",
			driver:   "Generic Text",
			printer:  "Line Matrix",
			wantType: "serial",
		},
		{
			name:     "serial LPT non-thermal",
			port:     "LPT1",
			driver:   "Generic Text",
			printer:  "Dot Matrix",
			wantType: "serial",
		},
		{
			name:     "serial word in port non-thermal",
			port:     "SerialPort0",
			driver:   "Industrial Scale",
			printer:  "Scale Indicator",
			wantType: "serial",
		},
		{
			name:     "unknown printer port and driver",
			port:     "BLUETOOTH_DEV1",
			driver:   "Custom Device Driver",
			printer:  "Unknown Peripheral",
			wantType: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectPrinterType(tt.port, tt.driver, tt.printer)
			if got != tt.wantType {
				t.Errorf("detectPrinterType(%q, %q, %q) = %q, want %q", tt.port, tt.driver, tt.printer, got, tt.wantType)
			}
		})
	}
}

func TestFilterPhysicalPrinters(t *testing.T) {
	t.Run("filters out virtual printers", func(t *testing.T) {
		input := []PrinterDetail{
			{Name: "Microsoft Print to PDF", IsVirtual: true, PrinterType: "virtual"},
			{Name: "Microsoft XPS Document Writer", IsVirtual: true, PrinterType: "virtual"},
			{Name: "OneNote", IsVirtual: true, PrinterType: "virtual"},
			{Name: "EPSON TM-T88V", IsVirtual: false, PrinterType: "thermal"},
			{Name: "HP LaserJet Pro", IsVirtual: false, PrinterType: "network"},
			{Name: "Canon TS3300", IsVirtual: false, PrinterType: "usb"},
		}

		got := FilterPhysicalPrinters(input)
		if len(got) != 3 {
			t.Fatalf("FilterPhysicalPrinters returned %d printers, want 3", len(got))
		}

		expectedNames := []string{"EPSON TM-T88V", "HP LaserJet Pro", "Canon TS3300"}
		for i, name := range expectedNames {
			if got[i].Name != name {
				t.Errorf("got[%d].Name = %q, want %q", i, got[i].Name, name)
			}
		}
	})

	t.Run("empty input returns empty slice", func(t *testing.T) {
		got := FilterPhysicalPrinters(nil)
		if len(got) != 0 {
			t.Errorf("FilterPhysicalPrinters(nil) returned %d items, want 0", len(got))
		}

		got = FilterPhysicalPrinters([]PrinterDetail{})
		if len(got) != 0 {
			t.Errorf("FilterPhysicalPrinters([]) returned %d items, want 0", len(got))
		}
	})
}

func TestFilterThermalPrinters(t *testing.T) {
	t.Run("filters only thermal printers", func(t *testing.T) {
		input := []PrinterDetail{
			{Name: "EPSON TM-T88V", PrinterType: "thermal"},
			{Name: "Star TSP100", PrinterType: "thermal"},
			{Name: "POS-80 Printer", PrinterType: "thermal"},
			{Name: "58mm Mini Receipt", PrinterType: "thermal"},
			{Name: "Microsoft Print to PDF", PrinterType: "virtual"},
			{Name: "HP LaserJet Pro", PrinterType: "network"},
			{Name: "Canon TS3300", PrinterType: "usb"},
			{Name: "Generic Serial", PrinterType: "serial"},
		}

		got := FilterThermalPrinters(input)
		if len(got) != 4 {
			t.Fatalf("FilterThermalPrinters returned %d printers, want 4", len(got))
		}

		expectedNames := []string{"EPSON TM-T88V", "Star TSP100", "POS-80 Printer", "58mm Mini Receipt"}
		for i, name := range expectedNames {
			if got[i].Name != name {
				t.Errorf("got[%d].Name = %q, want %q", i, got[i].Name, name)
			}
		}
	})

	t.Run("empty input returns empty slice", func(t *testing.T) {
		got := FilterThermalPrinters(nil)
		if len(got) != 0 {
			t.Errorf("FilterThermalPrinters(nil) returned %d items, want 0", len(got))
		}

		got = FilterThermalPrinters([]PrinterDetail{})
		if len(got) != 0 {
			t.Errorf("FilterThermalPrinters([]) returned %d items, want 0", len(got))
		}
	})
}

func TestFindThermalPrinters(t *testing.T) {
	printers, err := FindThermalPrinters()
	if runtime.GOOS != "windows" {
		if err == nil {
			t.Errorf("FindThermalPrinters() on %s expected error, got nil (printers: %v)", runtime.GOOS, printers)
		}
	} else {
		if err != nil {
			t.Logf("FindThermalPrinters() on Windows returned error: %v", err)
		} else {
			for _, p := range printers {
				if p.PrinterType != "thermal" {
					t.Errorf("FindThermalPrinters() returned non-thermal printer: %+v", p)
				}
			}
		}
	}
}
