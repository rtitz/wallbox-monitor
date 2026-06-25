FROM alpine:3.19

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# create directory and set ownership to appuser:appgroup
RUN mkdir -p /home/appuser/wallbox-monitor \
 && chown -R appuser:appgroup /home/appuser/wallbox-monitor \
 && chmod 755 /home/appuser/wallbox-monitor

# Copy binary and set permissions
COPY bin/wallbox-monitor_linux-arm64 /usr/local/bin/wallbox-monitor
RUN chown appuser:appgroup /usr/local/bin/wallbox-monitor && chmod 755 /usr/local/bin/wallbox-monitor

USER appuser

EXPOSE 2112

ENTRYPOINT ["/usr/local/bin/wallbox-monitor"]
