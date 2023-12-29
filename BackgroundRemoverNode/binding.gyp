{
    "targets": [
        {
            "target_name": "BackgroundRemover",
            "cflags!": ["-fno-exceptions"],
            "cflags_cc!": ["-fno-exceptions"],
            "sources": ["lib/native/Binding.cc"],
            "include_dirs": [
                "<!@(node -p \"require('node-addon-api').include\")",
            ],
            'defines': ['NAPI_DISABLE_CPP_EXCEPTIONS'],
            'libraries': [],
            'msvs_settings': {
                'VCCLCompilerTool': {
                    'AdditionalOptions': ['-std:c++20', ]
                },
            }
        }
    ]
}
