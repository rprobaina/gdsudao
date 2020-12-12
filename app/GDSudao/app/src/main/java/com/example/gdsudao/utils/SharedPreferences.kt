package com.example.gdsudao.utils

import android.app.Activity
import android.content.Context
import android.util.Log
import android.widget.Toast
import com.example.gdsudao.model.Area
import com.google.gson.Gson
import java.lang.Exception
import com.google.gson.reflect.TypeToken as TypeToken

class SharedPreferences {

    fun RemoverAreaLista(context: Context, area: Area){
        var listaAreas = RecuperarListaAreas(context)
        Log.println(Log.DEBUG, "SHARED", listaAreas.toString())
        listaAreas.remove(area)
        Log.println(Log.DEBUG, "SHARED", listaAreas.toString())

        var pref = context.getSharedPreferences("GDSUDAO_PREFERENCIAS", Context.MODE_PRIVATE)
        var prefEditor = pref.edit()
        var gson = Gson()
        var json = gson.toJson(listaAreas)

        prefEditor.putString("area", json)
        prefEditor.commit()
    }

    fun RemoverAllAreaLista(context: Context){
        var listaAreas = RecuperarListaAreas(context)

        var pref = context.getSharedPreferences("GDSUDAO_PREFERENCIAS", Context.MODE_PRIVATE)
        var prefEditor = pref.edit()

        var json = ""

        prefEditor.putString("area", json)
        prefEditor.commit()
    }

    fun SalvarAreaLista(context: Context, area: Area){
        var listaAreas = RecuperarListaAreas(context)
        Log.println(Log.DEBUG, "SHARED", listaAreas.toString())
        //listaAreas.add(Area("asasdasdasd", "23232"))
        listaAreas.add(area)
        Log.println(Log.DEBUG, "SHARED", listaAreas.toString())

        var pref = context.getSharedPreferences("GDSUDAO_PREFERENCIAS", Context.MODE_PRIVATE)
        var prefEditor = pref.edit()
        var gson = Gson()
        var json = gson.toJson(listaAreas)

        prefEditor.putString("area", json)
        prefEditor.commit()
    }

    fun RecuperarListaAreas(context: Context): ArrayList<Area> {
        var listaAreas : ArrayList<Area> = ArrayList()
        var pref = context.getSharedPreferences("GDSUDAO_PREFERENCIAS", Context.MODE_PRIVATE)
        var gson = Gson()
        var json = pref.getString("area", "")

        try {
            val tipo = object : TypeToken<List<Area>>() {}.type
            listaAreas = gson.fromJson(json, tipo)
            Log.println(Log.DEBUG, "SHARED", listaAreas.toString())
        }catch (e: Exception){

        }

        //Toast.makeText(context, json.toString(), Toast.LENGTH_SHORT).show()
        return listaAreas
    }
}