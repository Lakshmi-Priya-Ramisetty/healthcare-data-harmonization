module github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/tools/notebook/extensions/wstl/service

go 1.14

replace github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine => ../../../../../mapping_engine

replace github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/proto => ../../../../../mapping_engine/proto

replace github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/util => ../../../../../mapping_engine/util

replace github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/transform => ../../../../../mapping_engine/transform

replace github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_language => ../../../../../mapping_language

replace github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_language/transpiler => ../../../../../mapping_language/transpiler

replace github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/tools/notebook/extensions/wstl/proto => ../proto

require (
    cloud.google.com/go/storage v1.13.0
    github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/proto v0.0.0-20210205224744-b12bbd79968b
    github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/transform v0.0.0-20210205224744-b12bbd79968b
    github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/util v0.0.0-20210205224744-b12bbd79968b
    github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/tools/notebook/extensions/wstl/proto v0.0.0-20210205224744-b12bbd79968b
    github.com/google/fhir/go v0.0.0-20210120234235-b7cfb32dc82f
    github.com/google/go-cmp v0.5.4
    google.golang.org/genproto v0.0.0-20210203152818-3206188e46ba
    google.golang.org/grpc v1.35.0
    google.golang.org/protobuf v1.25.0
)
