# API Script Guide

## 1. Overview

This document defines the guidelines for writing API scripts in the Geelato platform.

> Important: Platform capabilities are continuously updated. Please configure MCP for real-time capabilities.

## 2. File Structure

api/
`└── {groupName}/`
    `├── {apiCode}.api.js`    # API script
    `└── {apiCode}.api.json`  # API metadata (optional)

## 3. Script Template

(function() {
    var param1 = $params.param1;

    // Business logic here

    return {
        code: 200,
        message: "success",
        data: {
            result: true
        }
    };
})();
