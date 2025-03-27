# FastZip4J
The FastZip4J is a library which uses the GoLang fastzip ``https://github.com/saracen/fastzip`` behind for performance and follow the same approach as ``fastzip``:
- Permissions, ownership (uid, gid on linux/unix) and modification times are preserved.
- Buffers used for copying files are recycled to reduce allocations.
- Files are archived and extracted concurrently.
- ``github.com/klauspost/compress/flate`` library is used for compression and decompression.

Note: FastZip4j does not support file overwriting. If you attempt to archive files with the same name, an error will be raised.

# Example
## Archive Single File
When archiving a single file, the FastZip4j will automatically create a ZIP file if one does not already exist. If the ZIP file is present, the new file will be added to the existing archive.
```java
/*
CompressionLevel must be between 1 (BestSpeed) and 9 (BestCompression). 
Higher levels typically run slower but compress more.
*/
var compressionLevel = 1;
FastZip4j.archiveFile(
  new File("path/to/your/file.txt"),
  new File("path/to/your/zipfile.zip"),
  compressionLevel);
```

## Archive Folder
When archiving folder, the FastZip4j will automatically create a ZIP file if one does not already exist. If the ZIP file is present, files in the folder will be added to the existing archive.
```java
/*
CompressionLevel must be between 1 (BestSpeed) and 9 (BestCompression). 
Higher levels typically run slower but compress more.
*/
var compressionLevel = 1;
FastZip4j.archiveDir(
  new File("path/to/your/folder/"),
  new File("path/to/your/zipfile.zip"),
  compressionLevel);
```

## Extractor
```java
FastZip4j.extract(
  new File("path/to/your/zipfile.zip"),
  new File("path/to/your/destination/folder/"));
```