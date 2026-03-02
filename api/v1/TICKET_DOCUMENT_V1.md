# POSTER: Document Format v1.0

## Estructura del Documento

El formato de documento para impresoras POS define una estructura JSON para describir trabajos de impresión complejos.

### Documento Principal

```json
{
  "version": "1.0",
  "profile": {
    /* ProfileConfig */
  },
  "debug_log": false,
  "commands": [
    /* Array de comandos */
  ]
}
```

| Campo       | Tipo          | Requerido | Descripción                                |
|-------------|---------------|-----------|--------------------------------------------|
| `version`   | string        | ✓         | Versión del formato (patrón: `^\d+\.\d+$`) |
| `profile`   | ProfileConfig | ✓         | Configuración del perfil de impresora      |
| `debug_log` | boolean       |           | Habilita logs de depuración                |
| `commands`  | Command[]     | ✓         | Lista de comandos a ejecutar (mínimo 1)    |

### ProfileConfig

Define las características de la impresora:

```json
{
  "model": "58mm PT-210",
  "paper_width": 58,
  "code_table": "PC850",
  "dpi": 203,
  "has_qr": true
}
```

| Campo         | Tipo    | Requerido | Descripción                      | Default | Valores                     |
|---------------|---------|-----------|----------------------------------|---------|-----------------------------|
| `model`       | string  | ✓         | Identificador del modelo         |         |                             |
| `paper_width` | integer |           | Ancho del papel en mm            | 80      | 58, 72, 80, 100, 112, 120   |
| `code_table`  | string  |           | Tabla de caracteres              | WPC1252 | PC437, PC850, WPC1252, etc. |
| `dpi`         | integer |           | Resolución en puntos por pulgada | 203     | 203, 300, 600               |
| `has_qr`      | boolean |           | Indica soporte nativo de QR      | false   |                             |

## Comandos Disponibles

### Command Structure

| Campo  | Tipo   | Requerido | Descripción                   | Valores                                               |
|--------|--------|-----------|-------------------------------|-------------------------------------------------------|
| `type` | string | ✓         | Tipo de comando               | text, image, separator, feed, cut, qr, table, barcode |
| `data` | object | ✓         | Datos específicos del comando | Varía según el tipo                                   |

### 1. Text Command

Imprime texto con formato opcional:

```json
{
  "type": "text",
  "data": {
    "content": {
      "text": "Texto principal",
      "content_style": {
        "bold": true,
        "size": "2x2"
      },
      "align": "center"
    },
    "label": {
      "text": "Etiqueta",
      "separator": ": ",
      "label_style": {
        "bold": false
      },
      "align": "left"
    },
    "new_line": true
  }
}
```

**TextCommand:**

| Campo      | Tipo    | Requerido | Descripción                      | Default |
|------------|---------|-----------|----------------------------------|---------|
| `content`  | Content | ✓         | Contenido principal del texto    |         |
| `label`    | Label   |           | Etiqueta opcional                |         |
| `new_line` | boolean |           | Salto de línea después del texto | true    |

**Content:**

| Campo           | Tipo      | Requerido | Descripción              | Default | Valores             |
|-----------------|-----------|-----------|--------------------------|---------|---------------------|
| `text`          | string    | ✓         | Texto a imprimir         |         |                     |
| `content_style` | TextStyle |           | Estilo del contenido     |         |                     |
| `align`         | string    |           | Alineación del contenido | left    | left, center, right |

**Label:**

| Campo         | Tipo      | Requerido | Descripción                          | Default | Valores             |
|---------------|-----------|-----------|--------------------------------------|---------|---------------------|
| `text`        | string    |           | Texto de la etiqueta                 | ""      |                     |
| `label_style` | TextStyle |           | Estilo de la etiqueta                |         |                     |
| `separator`   | string    |           | Separador entre etiqueta y contenido | ": "    |                     |
| `align`       | string    |           | Alineación de la etiqueta            | left    | left, center, right |

**TextStyle:**

| Campo       | Tipo    | Requerido | Descripción                    | Default | Valores                   |
|-------------|---------|-----------|--------------------------------|---------|---------------------------|
| `bold`      | boolean |           | Texto en negrita               | false   |                           |
| `size`      | string  |           | Tamaño del texto               | 1x1     | Patrón: `^([1-8]x[1-8])$` |
| `underline` | string  |           | Estilo de subrayado            | 0pt     | 0pt, 1pt, 2pt             |
| `inverse`   | boolean |           | Texto blanco sobre fondo negro | false   |                           |
| `font`      | string  |           | Fuente del texto               | A       | A, B                      |

### 2. Image Command

Imprime imágenes en formato base64:

```json
{
  "type": "image",
  "data": {
    "code": "base64_data...",
    "format": "png",
    "pixel_width": 384,
    "align": "center",
    "threshold": 128,
    "dithering": "atkinson",
    "scaling": "bilinear"
  }
}
```

| Campo         | Tipo    | Requerido | Descripción               | Default  | Valores             |
|---------------|---------|-----------|---------------------------|----------|---------------------|
| `code`        | string  | ✓         | Datos de imagen en Base64 |          |                     |
| `format`      | string  |           | Formato de imagen         |          | png, jpg, bmp       |
| `pixel_width` | integer |           | Ancho deseado en píxeles  | 128      | Mínimo: 1           |
| `align`       | string  |           | Alineación de la imagen   | center   | left, center, right |
| `threshold`   | integer |           | Umbral B/N (0-255)        | 128      | 0-255               |
| `dithering`   | string  |           | Algoritmo de dithering    | atkinson | threshold, atkinson |
| `scaling`     | string  |           | Algoritmo de escalado     | bilinear | bilinear, nns       |

### 3. Barcode Command

Genera códigos de barras:

```json
{
  "type": "barcode",
  "data": {
    "symbology": "CODE128",
    "data": "123456789",
    "width": 3,
    "height": 100,
    "hri_position": "below",
    "hri_font": "A",
    "align": "center"
  }
}
```

| Campo          | Tipo    | Requerido | Descripción           | Default | Valores                                                |
|----------------|---------|-----------|-----------------------|---------|--------------------------------------------------------|
| `symbology`    | string  | ✓         | Tipo de simbología    |         | upca, upce, ean13, ean8, code39, code128, itf, codabar |
| `data`         | string  | ✓         | Datos a codificar     |         | 1-25 caracteres                                        |
| `width`        | integer |           | Ancho del módulo      |         | 2-6                                                    |
| `height`       | integer |           | Altura en puntos      |         | 1-255                                                  |
| `hri_position` | string  |           | Posición del HRI      | below   | none, above, below, both                               |
| `hri_font`     | string  |           | Fuente para HRI       | A       | A, B                                                   |
| `align`        | string  |           | Alineación del código | center  | left, center, right                                    |

### 4. QR Command

Genera códigos QR:

```json
{
  "type": "qr",
  "data": {
    "data": "https://example.com",
    "human_text": "Visita nuestro sitio",
    "pixel_width": 200,
    "correction": "M",
    "align": "center",
    "logo": "base64_logo...",
    "circle_shape": false
  }
}
```

| Campo          | Tipo    | Requerido | Descripción                                    | Default | Valores             |
|----------------|---------|-----------|------------------------------------------------|---------|---------------------|
| `data`         | string  | ✓         | Datos del QR (URL, texto, etc.)                |         |                     |
| `human_text`   | string  |           | Texto a mostrar debajo del QR                  |         |                     |
| `pixel_width`  | integer |           | Ancho del QR en píxeles                        | 128     | Mínimo: 87          |
| `correction`   | string  |           | Nivel de corrección de errores                 | Q       | L, M, Q, H          |
| `align`        | string  |           | Alineación del QR                              | center  | left, center, right |
| `logo`         | string  |           | Logo en Base64                                 |         |                     |
| `circle_shape` | boolean |           | Usar bloques circulares (solo para QR < 256px) |         |                     |

### 5. Table Command

Crea tablas formateadas:

```json
{
  "type": "table",
  "data": {
    "definition": {
      "columns": [
        {
          "name": "Producto",
          "width": 20,
          "align": "left"
        },
        {
          "name": "Precio",
          "width": 10,
          "align": "right"
        }
      ]
    },
    "show_headers": true,
    "rows": [
      [
        "Café",
        "$3.50"
      ],
      [
        "Sandwich",
        "$8.00"
      ]
    ],
    "options": {
      "header_bold": true,
      "word_wrap": true,
      "column_spacing": 1,
      "align": "center",
      "auto_reduce": true
    }
  }
}
```

**TableCommand:**

| Campo          | Tipo            | Requerido | Descripción                                 | Default |
|----------------|-----------------|-----------|---------------------------------------------|---------|
| `definition`   | TableDefinition | ✓         | Definición de la estructura de la tabla     |         |
| `show_headers` | boolean         |           | Mostrar encabezados                         | true    |
| `rows`         | array[][]       | ✓         | Filas de datos (array de arrays de strings) |         |
| `options`      | TableOptions    |           | Opciones de renderizado                     |         |

**TableDefinition:**

| Campo         | Tipo          | Requerido | Descripción                                |
|---------------|---------------|-----------|--------------------------------------------|
| `columns`     | TableColumn[] | ✓         | Definición de columnas (mínimo 1)          |
| `paper_width` | integer       |           | Ancho del papel en caracteres (mínimo:  1) |

**TableColumn:**

| Campo   | Tipo    | Requerido | Descripción          | Default | Valores             |
|---------|---------|-----------|----------------------|---------|---------------------|
| `name`  | string  | ✓         | Nombre de la columna |         |                     |
| `width` | integer | ✓         | Ancho en caracteres  |         | Mínimo: 1           |
| `align` | string  |           | Alineación del texto | center  | left, center, right |

**TableOptions:**

| Campo            | Tipo    | Requerido | Descripción                                                     | Default | Valores             |
|------------------|---------|-----------|-----------------------------------------------------------------|---------|---------------------|
| `header_bold`    | boolean |           | Encabezados en negrita                                          | true    |                     |
| `word_wrap`      | boolean |           | Ajuste automático de texto                                      | true    |                     |
| `column_spacing` | integer |           | Espacios entre columnas                                         | 1       | Mínimo:  0          |
| `align`          | string  |           | Alineación de la tabla                                          | center  | left, center, right |
| `auto_reduce`    | boolean |           | Reducir automáticamente anchos de columna para ajustar al papel | true    |                     |

### 6. Separator Command

Crea líneas separadoras:

```json
{
  "type": "separator",
  "data": {
    "char": "-",
    "length": 48
  }
}
```

| Campo    | Tipo    | Requerido | Descripción                    | Default | Valores |
|----------|---------|-----------|--------------------------------|---------|---------|
| `char`   | string  |           | Carácter(es) para el separador | "- "    |         |
| `length` | integer |           | Longitud en caracteres         | 48      | 1-255   |

### 7. Feed Command

Avanza el papel:

```json
{
  "type": "feed",
  "data": {
    "lines": 3
  }
}
```

| Campo   | Tipo    | Requerido | Descripción                | Valores |
|---------|---------|-----------|----------------------------|---------|
| `lines` | integer | ✓         | Número de líneas a avanzar | 1-255   |

### 8. Cut Command

Corta el papel:

```json
{
  "type": "cut",
  "data": {
    "mode": "partial",
    "feed": 5
  }
}
```

| Campo  | Tipo    | Requerido | Descripción                      | Default | Valores       |
|--------|---------|-----------|----------------------------------|---------|---------------|
| `mode` | string  |           | Tipo de corte                    |         | full, partial |
| `feed` | integer |           | Líneas a avanzar antes del corte | 2       | 0-255         |

### 9. Raw Command

Envía bytes directamente a la impresora sin procesamiento:

⚠️ **ADVERTENCIA**: Comando avanzado. Uso incorrecto puede dañar la configuración de la impresora.

```json
{
  "type": "raw",
  "data": {
    "hex": "1B 40",
    "format": "hex",
    "comment": "Reset printer",
    "safe_mode": false
  }
}
```

| Campo       | Tipo    | Requerido | Descripción                         | Default | Valores     |
|-------------|---------|-----------|-------------------------------------|---------|-------------|
| `hex`       | string  | ✓         | Bytes en formato especificado       |         |             |
| `format`    | string  |           | Formato de entrada                  | hex     | hex, base64 |
| `comment`   | string  |           | Descripción del comando             |         |             |
| `safe_mode` | boolean |           | Habilitar validaciones de seguridad | false   |             |

**Formatos de Hex Soportados:**

- Espacios: `"1B 40"`
- Sin espacios: `"1B40"`
- Con comas: `"1B,40"`
- Con prefijo: `"0x1B 0x40"`

**Límites:**

- Máximo 4096 bytes por comando
- Solo caracteres hexadecimales válidos

### 10. Pulse Command

Genera un pulso eléctrico (apertura de cajón).

```json
{
  "type": "pulse",
  "data": {
    "pin": 0,
    "on_time": 50,
    "off_time": 100
  }
}

```

| Campo      | Tipo | Descripción              | Default |
|------------|------|--------------------------|---------|
| `pin`      | int  | Pin del conector (0 o 1) | 0       |
| `on_time`  | int  | Tiempo encendido en ms   | 50      |
| `off_time` | int  | Tiempo apagado en ms     | 100     |

### 11. Beep Command

Emite un sonido de alerta (si la impresora cuenta con buzzer).

```json
{
  "type": "beep",
  "data": {
    "times": 2,
    "lapse": 1
  }
}

```

| Campo   | Tipo | Descripción                  | Default |
|---------|------|------------------------------|---------|
| `times` | int  | Cantidad de pitidos          | 1       |
| `lapse` | int  | Factor de duración/intervalo | 1       |

## Ejemplo Completo

```json
{
  "version": "1.0",
  "profile": {
    "model": "58mm PT-210",
    "paper_width": 58,
    "code_table": "PC850",
    "dpi": 203,
    "has_qr": true
  },
  "commands": [
    {
      "type": "text",
      "data": {
        "content": {
          "text": "RECEIPT",
          "content_style": {
            "bold": true,
            "size": "2x2"
          },
          "align": "center"
        }
      }
    },
    {
      "type": "separator",
      "data": {
        "char": "=",
        "length": 32
      }
    },
    {
      "type": "table",
      "data": {
        "definition": {
          "columns": [
            {
              "name": "Item",
              "width": 20,
              "align": "left"
            },
            {
              "name": "Price",
              "width": 12,
              "align": "right"
            }
          ],
          "paper_width": 48
        },
        "show_headers": false,
        "rows": [
          [
            "Coffee",
            "$3.50"
          ],
          [
            "Muffin",
            "$4.25"
          ]
        ],
        "options": {
          "header_bold": true,
          "word_wrap": true,
          "column_spacing": 1
        }
      }
    },
    {
      "type": "barcode",
      "data": {
        "symbology": "CODE128",
        "data": "INV-2024-001",
        "width": 3,
        "height": 80,
        "hri_position": "below",
        "align": "center"
      }
    },
    {
      "type": "qr",
      "data": {
        "data": "https://example.com/receipt/12345",
        "human_text": "Scan for digital receipt",
        "pixel_width": 200,
        "correction": "M",
        "align": "center"
      }
    },
    {
      "type": "feed",
      "data": {
        "lines": 3
      }
    },
    {
      "type": "cut",
      "data": {
        "mode": "partial",
        "feed": 2
      }
    }
  ]
}
```

## Notas de Validación

- **Versión**: Debe seguir el patrón `^\d+\.\d+$` (ej: "1.0", "2.1")
- **ProfileConfig.model**: Es el único campo requerido en el perfil
- **Commands**: Debe contener al menos un comando
- **Barcode.data**: Limitado a 1-25 caracteres según el schema
- **QR.pixel_width**: Mínimo 87 píxeles
- **QR.circle_shape**: Solo recomendado para códigos QR mayores a 256px de ancho
- **Table.columns**: Debe tener al menos una columna definida
- **Valores por defecto**: Se aplican cuando el campo no está presente en el JSON