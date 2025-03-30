package com.fastzip4j;

import com.sun.jna.Native;

import java.nio.file.Path;

/**
 * The {@code FastZip4j} class provides utility methods for archiving files and directories
 * into ZIP files and extracting ZIP files to a specified directory. It uses native library
 * bindings for fast and efficient compression and decompression.
 */
public class FastZip4j {
    private static final String LIBRARY_NAME = "fastzip4j";
    private static final FastZip4jLib fastzip4jLib = Native.load(LIBRARY_NAME, FastZip4jLib.class);

    private FastZip4j() {
    }

    /**
     * Archives a single file into a ZIP file with the specified compression level.
     *
     * @param sourceFile       the file to be archived
     * @param zipFile          the destination ZIP file
     * @param compressionLevel the compression level. CompressionLevel must be between 1 (BestSpeed) and 9 (BestCompression). Higher levels typically run slower but compress more.
     */
    public static void archiveFile(Path sourceFile, Path zipFile, int compressionLevel) {
        archiveFile(sourceFile.toAbsolutePath().normalize().toString(), zipFile.toAbsolutePath().normalize().toString(), compressionLevel);
    }

    /**
     * Archives a directory and its contents into a ZIP file with the specified compression level.
     *
     * @param sourceDir        the directory to be archived
     * @param zipFile          the destination ZIP file
     * @param compressionLevel the compression level. CompressionLevel must be between 1 (BestSpeed) and 9 (BestCompression). Higher levels typically run slower but compress more.
     */
    public static void archiveDir(Path sourceDir, Path zipFile, int compressionLevel) {
        archiveDir(sourceDir.toAbsolutePath().normalize().toString(), zipFile.toAbsolutePath().normalize().toString(), compressionLevel);
    }

    /**
     * Extracts the contents of a ZIP file into the specified destination directory.
     *
     * @param zipFile              the ZIP file to be extracted
     * @param destinationDirectory the destination directory where the contents will be extracted
     */
    public static void extract(Path zipFile, Path destinationDirectory) {
        extract(zipFile.toAbsolutePath().normalize().toString(), destinationDirectory.toAbsolutePath().normalize().toString());
    }

    private static void archiveFile(String sourceFile, String zipFile, int compressionLevel) {
        fastzip4jLib.ArchiveFile(sourceFile, zipFile, compressionLevel);
    }

    private static void archiveDir(String sourceDir, String zipFile, int compressionLevel) {
        fastzip4jLib.ArchiveDir(sourceDir, zipFile, compressionLevel);
    }

    private static void extract(String zipFile, String destinationDirectory) {
        fastzip4jLib.Extract(zipFile, destinationDirectory);
    }
}
