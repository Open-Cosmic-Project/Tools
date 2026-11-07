#pragma once

#include <string>
#include <fstream>
#include <mutex>
#include <chrono>
#include <iomanip>

namespace logengine {

enum class Level { DEBUG, INFO, WARN, ERROR };

class Logger {
public:
    Logger(const std::string& file_path);
    ~Logger();

    void log(Level level, const std::string& message);

    // Convenience helpers
    void debug(const std::string& msg) { log(Level::DEBUG, msg); }
    void info(const std::string& msg)  { log(Level::INFO,  msg); }
    void warn(const std::string& msg)  { log(Level::WARN,  msg); }
    void error(const std::string& msg) { log(Level::ERROR, msg); }

private:
    std::ofstream out_;        // Stream to log file (or stdout)
    std::mutex mtx_;           // Protect concurrent writes
    std::string level_to_string(Level lvl) const;
    std::string timestamp() const;
};

} // namespace logengine
