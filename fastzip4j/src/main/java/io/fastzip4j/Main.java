package io.fastzip4j;

import java.io.File;

public class Main {
    public static void main(String[] args) {
        FastZip4j.extract(
                new File("path/to/your/file"),
                new File("path/to/your/zipfile.zip"));
    }
}
