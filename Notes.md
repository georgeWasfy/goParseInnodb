## 1. FUNDAMENTAL CONCEPTS

### 1.1 Spaces (Tablespaces)

**Definition:** A space is InnoDB's logical file unit.

- **Physical:** Can be multiple OS files (ibdata1, ibdata2) but treated as one logical file
- **Identification:** 32-bit integer space ID
- **System Space:** Always space ID = 0 (ibdata1)
- **Per-Table Spaces:** Each .ibd file is a complete space with one table


### 1.2 Pages

**Definition:** The fundamental unit of storage in InnoDB.

**SPECIFICATIONS:**
- **Size:** 16 KiB (16,384 bytes)
- **Addressing:** 32-bit page number (offset from start of space)
- **Calculation:** Page N is at file offset: `N × 16384`
- **Space Limit:** 2^32 pages × 16 KiB = 64 TiB per space

**Example Calculations:**
```
Page 0: Offset 0
Page 1: Offset 16384 (0x4000)
Page 2: Offset 32768 (0x8000)
Page 3: Offset 49152 (0xC000)
```

## 2. PAGE STRUCTURE

Every single page in InnoDB has this EXACT structure:

```
┌──────────────────────────────┐
│   FIL Header (38 bytes)      │  Offset 0
├──────────────────────────────┤
│                              │
│   Page Body                  │  Offset 38
│   (structure varies by type) │
│                              │
├──────────────────────────────┤
│   FIL Trailer (8 bytes)      │  Offset 16376
└──────────────────────────────┘
```
### 2.1 FIL Header (38 bytes) - EXACT STRUCTURE

**Byte Layout:**
```
Offset | Size | Field Name           | Description
-------|------|---------------------|------------------------------------
0      | 4    | Checksum            | 32-bit checksum of page
4      | 4    | Page Number         | Page offset (confirms position)
8      | 4    | Previous Page       | Prev page number in linked list
12     | 4    | Next Page           | Next page number in linked list
16     | 8    | LSN                 | Log Sequence Number (last mod)
24     | 2    | Page Type           | Type determines page body structure
26     | 8    | Flush LSN           | Only used in page 0 of space 0
34     | 4    | Space ID            | Which space this page belongs to
```

### 2.2 Page Types

**Page Types for Implementation:**
```go
const (
    PageTypeAllocated = 0      // Freshly allocated, not yet used
    PageTypeFSPHDR   = 8       // File space header (page 0)
    PageTypeXDES     = 8       // Extent descriptor (same structure as FSP_HDR)
    PageTypeIBUFBitmap = 5     // Insert buffer bitmap
    PageTypeINODE    = 3       // Index node (file segment inode)
    PageTypeINDEX    = 17855   // B-tree node (0x45BF) - YOUR DATA IS HERE
    PageTypeTRXSYS   = 7       // Transaction system header
    PageTypeSYS      = 6       // System page
    PageTypeBLOB     = 10      // Uncompressed BLOB page
)
```
### 2.3 FIL Trailer (8 bytes) - EXACT STRUCTURE

Located at offset 16376 (last 8 bytes of page):

```
Offset | Size | Field Name      | Description
-------|------|----------------|--------------------------------
16376  | 4    | Old Checksum   | Deprecated checksum (ignore)
16380  | 4    | LSN Low 32-bit | Low 32-bits of LSN from header
```

## 3. SPACE STRUCTURE

### 3.1 First Three Pages (ALL spaces have these)

**RULE:** Every InnoDB space file MUST start with:

```
Page 0: FSP_HDR    (File Space Header)
Page 1: IBUF_BITMAP (Insert Buffer Bitmap)  
Page 2: INODE      (Index Node Page)
```

**For .ibd file (per-table space):**
```
Page 0: FSP_HDR       - Space metadata
Page 1: IBUF_BITMAP   - Insert buffer bookkeeping
Page 2: INODE         - File segment information
Page 3: INDEX (ROOT)  - PRIMARY KEY root page ← YOUR DATA STARTS HERE
Page 4: INDEX (ROOT)  - First secondary index root (if exists)
Page 5+: INDEX (LEAF) - Actual data pages
```

### 3.2 FSP_HDR Structure (Page 0) Needs further in-depth viewing

**Fields in FSP_HDR:**
- Space size (number of pages)
- Lists of free extents
- Lists of fragmented extents
- Lists of full extents
- Extent descriptors for first 256 extents (16,384 pages)

### 3.3 XDES Pages (Extent Descriptor) Needs further in-depth viewing

**Rule:** Every 16,384 pages need an XDES page
**Note**: The structure of XDES and FSP_HDR pages is identical, except that the FSP header structure is zeroed-out in XDES pages

**Locations:**
```
Page 0:     FSP_HDR (contains XDES for extents 0-255)
Page 16384: XDES (extents 256-511)
Page 32768: XDES (extents 512-767)
...
```
### 3.4 INODE Page (Page 2) Needs further in-depth viewing

**Purpose:** Stores file segment (inode) entries

**What's a File Segment?**
- Logical grouping of extents and fragment pages
- Each index has TWO inodes: one for leaf pages, one for non-leaf
- 85 inode entries per INODE page

## 4. THE SYSTEM SPACE STRUCTURE (ibdata1)
**Skipped for now, I dont think I need to parse this**
```
Page 0:       FSP_HDR
Page 1:       IBUF_BITMAP
Page 2:       INODE
Page 3:       SYS (Insert buffer header)
Page 4:       INDEX (Insert buffer tree)
Page 5:       TRX_SYS (Transaction system)
Page 6:       SYS (Rollback segment)
Page 7:       SYS (Data dictionary)
Pages 64-127: Double write buffer
Pages 128-191: Double write buffer
Page 192+:    Various indexes and data
```

## 5. PER-TABLE SPACE (.ibd) STRUCTURE

**This is OUR target for table `t`:**

```
┌─────────────────────────────────────────┐
│ Page 0: FSP_HDR                         │
├─────────────────────────────────────────┤
│ Page 1: IBUF_BITMAP                     │
├─────────────────────────────────────────┤
│ Page 2: INODE                           │
├─────────────────────────────────────────┤
│ Page 3: INDEX (Root of PRIMARY KEY)     │ 
│         - If table small: contains data │
│         - If table large: contains      │
│           pointers to child pages       │
├─────────────────────────────────────────┤
│ Page 4+: INDEX (Leaf pages)             │ ← DATA
│          - Actual row data              │
│          - Linked list (prev/next)      │
└─────────────────────────────────────────┘
```