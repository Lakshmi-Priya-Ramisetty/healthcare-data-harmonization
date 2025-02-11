module github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_language/test

go 1.14

replace github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_language => ../

replace github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine => ../../mapping_engine

replace github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_language/transpiler => ../transpiler

replace github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/transform => ../../mapping_engine/transform

replace github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/proto => ../../mapping_engine/proto

replace github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/util => ../../mapping_engine/util

require (
    github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/proto v0.0.0-00010101000000-000000000000
    github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/transform v0.0.0-00010101000000-000000000000
    github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_engine/util v0.0.0-00010101000000-000000000000
    github.com/Lakshmi-Priya-Ramisetty/healthcare-data-harmonization/mapping_language/transpiler v0.0.0-00010101000000-000000000000
    github.com/google/go-cmp v0.5.4
    google.golang.org/protobuf v1.25.0
)
