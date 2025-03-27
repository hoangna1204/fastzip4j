package io.fastzip4j;

import com.sun.jna.Native;

import java.io.File;

public class FastZip4j {
    private static final String LIBRARY_NAME = "fastzip4j";
    private static final FastZip4jLib fastzip4jLib = Native.load(LIBRARY_NAME, FastZip4jLib.class);

    private FastZip4j() {
    }

    public static void archiveFile(File sourceFile, File zipFile, int compressionLevel) {
        archiveFile(sourceFile.getAbsolutePath(), zipFile.getAbsolutePath(), compressionLevel);
    }

    public static void archiveFile(String sourceFile, String zipFile, int compressionLevel) {
        fastzip4jLib.ArchiveFile(sourceFile, zipFile, compressionLevel);
    }

    public static void archiveDir(File sourceDir, File zipFile, int compressionLevel) {
        archiveDir(sourceDir.getAbsolutePath(), zipFile.getAbsolutePath(), compressionLevel);
    }

    public static void archiveDir(String sourceDir, String zipFile, int compressionLevel) {
        fastzip4jLib.ArchiveDir(sourceDir, zipFile, compressionLevel);
    }

    public static void extract(File zipFile, File destinationDirectory) {
        extract(zipFile.getAbsolutePath(), destinationDirectory.getAbsolutePath());
    }

    public static void extract(String zipFile, String destinationDirectory) {
        fastzip4jLib.Extract(zipFile, destinationDirectory);
    }
}
