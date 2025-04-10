# FastZip4J
The FastZip4J is a library which is powered by GoLang behind for performance. The library is opinionated with a focus on speed:
- Permissions, ownership (uid, gid on linux/unix) and modification times are preserved.
- Buffers used for copying files are recycled to reduce allocations.
- Files are archived and extracted concurrently.
- File overwriting is not supported. If you attempt to archive files with the same name, an error will be raised.

# Example
## Import dependency

If you are using Maven
```xml
<dependency>
    <groupId>io.github.hoangna1204</groupId>
    <artifactId>fastzip4j</artifactId>
    <version>1.0.0</version>
</dependency>
```

If you are using Gradle
```groovy
implementation 'io.github.hoangna1204:fastzip4j:1.0.0'
```

## Archive Single File
When archiving a single file, the FastZip4j will automatically create a ZIP file if one does not already exist. If the ZIP file is present, the file will be added to the existing archive.
```java
/*
CompressionLevel must be between 1 (BestSpeed) and 9 (BestCompression). 
Higher levels typically run slower but compress more.
*/
var compressionLevel = 1;
FastZip4j.archiveFile(
  Path.of("path/to/your/file.txt"),
  Path.of("path/to/your/zipfile.zip"),
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
  Path.of("path/to/your/folder/"),
  Path.of("path/to/your/zipfile.zip"),
  compressionLevel);
```

## Extractor
```java
FastZip4j.extract(
  Path.of("path/to/your/zipfile.zip"),
  Path.of("path/to/your/destination/folder/"));
```