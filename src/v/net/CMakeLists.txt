v_cc_library(
  NAME net
  SRCS
    batched_output_stream.cc
    dns.cc
    transport.cc
    connection.cc
    server.cc
    probes.cc
    tls.cc
  DEPS
    Seastar::seastar
  )
