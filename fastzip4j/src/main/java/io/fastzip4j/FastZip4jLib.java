package io.fastzip4j;

import com.sun.jna.Library;

interface FastZip4jLib extends Library {
    void ArchiveFile(String sourceFile, String zipDestination, int compressionLevel);

    void ArchiveDir(String sourceDir, String zipFile, int compressionLevel);

    void Extract(String zipFile, String destinationDirectory);
}
