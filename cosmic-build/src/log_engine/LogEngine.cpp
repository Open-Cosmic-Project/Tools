#include "LogEngine.h"
#include <iostream>

namespace logengine {

Logger::Logger(const std::string& file_path) {
    if (file_path == "stdout" || file_path.empty()) {
        // Redirigir al stdout
        out_.basic_ios<char>::rdbuf(std::cout.rdbuf());
    } else {
        out_.open(file_path, std::ios::out | std::ios::app);
        if (!out_) {
            std::cerr << "logengine: no se pudo abrir " << file_path << ", usando stdout\n";
            out_.basic_ios<char>::rdbuf(std::cout.rdbuf());
        }
    }
}

Logger::~Logger() {
    if (out_.is_open()) {
        out_.close();
    }
}

void Logger::log(Level level, const std::string& message) {
    std::lock_guard<std::mutex> lock(mtx_);
    out_ << timestamp() << " [" << level_to_string(level) << "] " << message << "\n";
    out_.flush();
}

std::string Logger::level_to_string(Level lvl) const {
    switch (lvl) {
        case Level::DEBUG: return "DEBUG";
        case Level::INFO:  return "INFO";
        case Level::WARN:  return "WARN";
        case Level::ERROR: return "ERROR";
    }
    return "UNKNOWN";
}

std::string Logger::timestamp() const {
    using namespace std::chrono;
    auto now = system_clock::now();
    auto itt = system_clock::to_time_t(now);
    std::tm tm;
    localtime_r(&itt, &tm);
    auto ms = duration_cast<milliseconds>(now.time_since_epoch()) % 1000;
    std::ostringstream oss;
    oss << std::put_time(&tm, "%Y-%m-%d %H:%M:%S")
        << '.' << std::setfill('0') << std::setw(3) << ms.count();
    return oss.str();
}

} // namespace logengine

// --- Punto de entrada: lee líneas de stdin y las registra ---
int main(int argc, char* argv[]) {
    std::string log_path = "stdout";
    if (argc > 1) {
        log_path = argv[1];
    }

    logengine::Logger logger(log_path);
    logger.info("Motor de logs CosmicOS iniciado");

    std::string line;
    while (std::getline(std::cin, line)) {
        if (line.empty()) continue;
        // Formato: NIVEL mensaje  (ej: "INFO compilando cinitd")
        // Si no tiene nivel, se usa INFO por defecto
        if (line.substr(0, 5) == "DEBUG") {
            logger.debug(line.substr(6));
        } else if (line.substr(0, 4) == "WARN") {
            logger.warn(line.substr(5));
        } else if (line.substr(0, 5) == "ERROR") {
            logger.error(line.substr(6));
        } else if (line.substr(0, 4) == "INFO") {
            logger.info(line.substr(5));
        } else {
            logger.info(line);
        }
    }

    logger.info("Motor de logs CosmicOS finalizado");
    return 0;
}
